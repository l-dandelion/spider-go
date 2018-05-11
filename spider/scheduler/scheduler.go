package scheduler

import (
	"context"
	"fmt"
	"sync"

	"errors"
	"github.com/l-dandelion/spider-go/lib/library/buffer"
	"github.com/l-dandelion/spider-go/lib/library/cmap"
	"github.com/l-dandelion/spider-go/lib/library/parseurl"
	"github.com/l-dandelion/spider-go/spider/model/parsers"
	"github.com/l-dandelion/spider-go/spider/model/processors"
	"github.com/l-dandelion/spider-go/spider/module"
	"github.com/l-dandelion/spider-go/spider/module/data"
	"github.com/l-dandelion/spider-go/spider/module/stub"
	"github.com/l-dandelion/spider-go/spider/spider"
	"github.com/l-dandelion/yi-ants-go/lib/constant"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Scheduler interface {
	SchedulerName() string
	Start(initialReqs []*data.Request) error
	Pause() error
	Recover() error
	Stop() error
	Status() int8
	ErrorChan() <-chan error // get error
	Summary() Summary                    // get schduler summary
	SendReq(req *data.Request) bool
	SetDistributeQueue(pool buffer.Pool)
	SignRequest(request *data.Request)
	HasRequest(request *data.Request) bool
}

/*
 * create an instance of interface Scheduler by name
 */
func New(name string) Scheduler {
	bufferPool := buffer.NewPool(50, 10)
	return &myScheduler{
		name:             name,
		downloaderCounts: stub.NewModuleInternal(),
		errorBufferPool:  bufferPool,
	}
}

/*
 * create an instance of interface Scheduler
 */
func NewScheduler() Scheduler {
	return &myScheduler{}
}

/*
 * implementation of interface Scheduler
 */
type myScheduler struct {
	name              string
	maxDepth          uint32             // the max crawl depth
	acceptedDomainMap cmap.ConcurrentMap // accepted domain
	errorBufferPool   buffer.Pool        // error buffer pool
	urlMap            cmap.ConcurrentMap // url map
	ctx               context.Context    // used for stoping
	cancelFunc        context.CancelFunc // used for stoping
	status            int8               // running status
	statusLock        sync.RWMutex       // status lock
	distributeQeueu   buffer.Pool
	reportQueue       buffer.Pool
	initError         error
	respParsers       []module.ParseResponse
	itemProcessors    []module.ProcessItem
	downloaderCounts  stub.ModuleInternal
	initialReqs       []*data.Request
	beginAt           time.Time
	endAt             time.Time
}

/*
 * get scheduler name
 */
func (sched *myScheduler) SchedulerName() string {
	return sched.name
}

/*
 * initialize schduler
 */
func (sched *myScheduler) Init(job *spider.Job) (err error) {
	//check status
	log.Info("Check status for initialization...")
	oldStatus, err := sched.checkAndSetStatus(constant.RUNNING_STATUS_PREPARING)
	if err != nil {
		return
	}
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.initError = err
			sched.status = oldStatus
		} else {
			sched.status = constant.RUNNING_STATUS_PREPARED
		}
		sched.statusLock.Unlock()
	}()

	//check arguments
	log.Info("Check accepted domains...")
	if len(job.AcceptedDomains) == 0 {
		log.Info(ErrEmptyAcceptedDomainList)
		return ErrEmptyAcceptedDomainList
	}
	log.Info("Accepted domains are valid.")

	log.Info("Generate resposne parsers ...")
	sched.respParsers, err = parsers.GenParsersByModels(job.ParserModels)
	log.Info("Genetare response parsers success.")

	log.Info("Generate item processors")
	sched.itemProcessors, err = processors.GenProcessorsByModels(job.ProcessorModels)
	log.Info("Genetare item processors success.")

	// initialize internal fields
	log.Info("Initialize Scheduler's fields...")
	sched.maxDepth = job.MaxDepth
	log.Infof("-- Max depth: %d", sched.maxDepth)

	log.Info("Initialize requests...")
	initialUrls := parseurl.ParseReqUrl(job.InitialUrls, nil)
	sched.initialReqs = []*data.Request{}
	for _, urlStr := range initialUrls {
		httpReq, err := http.NewRequest("Get", urlStr, nil)
		if err != nil {
			return
		}
		req := data.NewRequest(httpReq)
		sched.initialReqs = append(sched.initialReqs, req)
	}
	log.Info("Initialize requests success.")

	sched.acceptedDomainMap, _ = cmap.NewConcurrentMap(1, nil)
	for _, domain := range job.AcceptedDomains {
		sched.acceptedDomainMap.Put(domain, struct{}{})
	}
	log.Infof("-- Accepted primay domains: %v", job.AcceptedDomains)

	sched.urlMap, _ = cmap.NewConcurrentMap(16, nil)
	log.Infof("-- URL map: length: %d, concurrency: %d", sched.urlMap.Len(), sched.urlMap.Concurrency())

	sched.resetContext()

	log.Info("Scheduler has been initialized.")
	return
}

/*
 * start scheduler
 */
func (sched *myScheduler) Start(isFirst bool) (err error) {
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Fatal Scheduler error: %s", p)
			log.Fatal(errMsg)
			err = errors.New(errMsg)
		}
	}()
	log.Info("Start Scheduler ...")
	log.Info("Check status for start ...")
	var oldStatus int8
	oldStatus, err = sched.checkAndSetStatus(constant.RUNNING_STATUS_STARTING)
	if err != nil {
		return
	}
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = constant.RUNNING_STATUS_STARTED
		}
		sched.statusLock.Unlock()
	}()
	log.Info("Get the primary domain...")

	for _, req := range sched.initialReqs {
		httpReq := req.HttpReq
		log.Infof("-- Host: %s", httpReq.Host)
		var primaryDomain string
		primaryDomain, err = getPrimaryDomain(httpReq.Host)
		if err != nil {
			return
		}
		ok, _ := sched.acceptedDomainMap.Put(primaryDomain, struct{}{})
		if ok {
			log.Infof("-- Primary domain: %s", primaryDomain)
		}
	}

	sched.download()
	sched.analyze()
	sched.pick()
	log.Info("The Scheduler has been started.")
	for _, req := range sched.initialReqs {
		sched.sendReq(req)
	}
	return nil
}

/*
 * pause scheduler
 */
func (sched *myScheduler) Pause() (err error) {
	//check status
	log.Info("Pause Scheduler ...")
	log.Info("Check status for pause ...")
	var oldStatus int8
	oldStatus, err = sched.checkAndSetStatus(constant.RUNNING_STATUS_PAUSING)
	if err != nil {
		return
	}
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = constant.RUNNING_STATUS_PAUSED
		}
		sched.statusLock.Unlock()
	}()
	log.Info("Scheduler has been paused.")
	return nil
}

/*
 * recover scheduler
 */
func (sched *myScheduler) Recover() (err error) {
	log.Info("Recover Scheduler ...")
	log.Info("Check status for recover ...")
	var oldStatus int8
	oldStatus, err = sched.checkAndSetStatus(constant.RUNNING_STATUS_STARTING)
	if err != nil {
		return
	}
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = constant.RUNNING_STATUS_STARTED
		}
		sched.statusLock.Unlock()
	}()
	log.Info("Scheduler has been recovered.")
	return nil
}

/*
 * stop scheduler
 */
func (sched *myScheduler) Stop() (err error) {
	log.Info("Stop Scheduler ...")
	log.Info("Check status for stop ...")
	var oldStatus int8
	oldStatus, err = sched.checkAndSetStatus(constant.RUNNING_STATUS_STOPPING)
	if err != nil {
		return
	}
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = constant.RUNNING_STATUS_STOPPED
		}
		sched.statusLock.Unlock()
	}()

	sched.cancelFunc()
	sched.errorBufferPool.Close()
	log.Info("Scheduler has been stopped.")
	return nil
}

/*
 * get error chan
 */
func (sched *myScheduler) ErrorChan() <-chan error {
	errBuffer := sched.errorBufferPool
	errCh := make(chan error, errBuffer.BufferCap())
	go func(errBuffer buffer.Pool, errCh chan error) {
		for {
			//stopped
			if sched.canceled() {
				close(errCh)
				break
			}
			//paused
			if sched.Status() == constant.RUNNING_STATUS_PAUSED {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			datum, err := errBuffer.Get()
			if err != nil {
				log.Warnln("The error buffer pool was closed. Break error reception.")
				close(errCh)
				break
			}
			err, ok := datum.(error)
			if !ok {
				err := fmt.Errorf("Incorrect error type: %T", datum)
				sched.sendError(err)
				continue
			}
			if sched.canceled() {
				close(errCh)
				break
			}
			errCh <- err
		}
	}(errBuffer, errCh)
	return errCh
}

/*
 * check status
 * set new status if check status success
 * return the old status if success, or an error return
 */
func (sched *myScheduler) checkAndSetStatus(wantedStatus int8) (oldStatus int8, err error) {
	sched.statusLock.Lock()
	defer sched.statusLock.Unlock()
	oldStatus = sched.status
	err = checkStatus(oldStatus, wantedStatus)
	if err == nil {
		sched.status = wantedStatus
	}
	return
}

/*
 * reset context
 */
func (sched *myScheduler) resetContext() {
	sched.ctx, sched.cancelFunc = context.WithCancel(context.Background())
}

/*
 * check whether the scheduler is stopped
 */
func (sched *myScheduler) canceled() bool {
	select {
	case <-sched.ctx.Done():
		return true
	default:
		return false
	}
}

/*
 * get running status
 */
func (sched *myScheduler) Status() int8 {
	sched.statusLock.RLock()
	defer sched.statusLock.RUnlock()
	return sched.status
}

/*
 * set distribute queue
 */
func (sched *myScheduler) SetDistributeQueue(pool buffer.Pool) {
	sched.distributeQeueu = pool
}

/*
 * sign request
 */
func (sched *myScheduler) SignRequest(req *data.Request) {
	sched.urlMap.Put(req.HttpReq.URL.String(), struct{}{})
}

/*
 * check whether it has request
 */
func (sched *myScheduler) HasRequest(req *data.Request) bool {
	return sched.urlMap.Get(req.HttpReq.URL.String()) != nil
}
