package queues

import (
	"github.com/l-dandelion/spider-go/lib/library/buffer"
	"github.com/l-dandelion/spider-go/spider/module/data"
	"sync"
)

var (
	DefaultBufferCap       = uint32(1000)
	DefaultMaxBufferNumber = uint32(100000)
	DefaultQueues = NewQueues()
)

type ContextQueue map[string]*data.Context

type Queues struct {
	sync.RWMutex
	SpiderRequestMap    map[string]buffer.Pool
	SpiderResponseMap   map[string]buffer.Pool
	SpiderItemMap       map[string]buffer.Pool
	NodeRequestMap      map[string]ContextQueue
	SpiderCrawlingQueue map[string]ContextQueue
}

func NewQueues() *Queues {
	return &Queues{
		SpiderRequestMap:    make(map[string]buffer.Pool),
		SpiderResponseMap:   make(map[string]buffer.Pool),
		SpiderItemMap:       make(map[string]buffer.Pool),
		NodeRequestMap:      make(map[string]ContextQueue),
		SpiderCrawlingQueue: make(map[string]ContextQueue),
	}
}

func (queue *Queues) PushRequest(ctx *data.Context) error {
	queue.Lock()
	spiderRequestQueue, ok := queue.SpiderRequestMap[ctx.SpiderName]
	if !ok {
		spiderRequestQueue, _ = buffer.NewPool(DefaultBufferCap, DefaultMaxBufferNumber)
		queue.SpiderRequestMap[ctx.SpiderName] = spiderRequestQueue
	}
	nodeReuqestQueue, ok := queue.NodeRequestMap[ctx.BackUpNodeName]
	if !ok {
		nodeReuqestQueue = make(ContextQueue)
		queue.NodeRequestMap[ctx.NodeName] = nodeReuqestQueue
	}
	queue.Unlock()

	err := spiderRequestQueue.Put(ctx)
	if err != nil {
		return err
	}
	nodeReuqestQueue[ctx.Unique()] = ctx
	return nil
}

func (queue *Queues) GetRequest(spiderName string) (*data.Context, error) {
	queue.Lock()
	spiderRequestQueue, ok := queue.SpiderRequestMap[spiderName]
	if !ok {
		spiderRequestQueue, _ = buffer.NewPool(DefaultBufferCap, DefaultMaxBufferNumber)
		queue.SpiderRequestMap[spiderName] = spiderRequestQueue
	}
	crawlingQueue, ok := queue.SpiderCrawlingQueue[spiderName]
	if !ok {
		crawlingQueue = make(ContextQueue)
		queue.SpiderCrawlingQueue[spiderName] = crawlingQueue
	}
	queue.Unlock()

	datum, err := spiderRequestQueue.Get()
	if err != nil {
		return nil, err
	}

	ctx, ok := datum.(*data.Context)
	if !ok {
		return nil, nil
	}
	crawlingQueue[ctx.Unique()] = ctx

	return ctx, nil
}

func (queue *Queues) PushResponse(ctx *data.Context) error {
	queue.Lock()
	spiderResponseQueue, ok := queue.SpiderResponseMap[ctx.SpiderName]
	if !ok {
		spiderResponseQueue, _ = buffer.NewPool(DefaultBufferCap, DefaultMaxBufferNumber)
		queue.SpiderResponseMap[ctx.SpiderName] = spiderResponseQueue
	}
	queue.Unlock()

	err := spiderResponseQueue.Put(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (queue *Queues) GetResponse(spiderName string) (*data.Context, error) {
	queue.Lock()
	spiderResponseQueue, ok := queue.SpiderResponseMap[spiderName]
	if !ok {
		spiderResponseQueue, _ = buffer.NewPool(DefaultBufferCap, DefaultMaxBufferNumber)
		queue.SpiderResponseMap[spiderName] = spiderResponseQueue
	}
	queue.Unlock()

	datum, err := spiderResponseQueue.Get()
	if err != nil {
		return nil, err
	}

	ctx, ok := datum.(*data.Context)
	if !ok {
		return nil, nil
	}

	return ctx, nil
}

func (queue *Queues) PushItem(ctx *data.Context) error {
	queue.Lock()
	spiderItemQueue, ok := queue.SpiderItemMap[ctx.SpiderName]
	if !ok {
		spiderItemQueue, _ = buffer.NewPool(DefaultBufferCap, DefaultMaxBufferNumber)
		queue.SpiderItemMap[ctx.SpiderName] = spiderItemQueue
	}
	queue.Unlock()

	err := spiderItemQueue.Put(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (queue *Queues) GetItem(spiderName string) (*data.Context, error) {
	queue.Lock()
	spiderItemQueue, ok := queue.SpiderItemMap[spiderName]
	if !ok {
		spiderItemQueue, _ = buffer.NewPool(DefaultBufferCap, DefaultMaxBufferNumber)
		queue.SpiderItemMap[spiderName] = spiderItemQueue
	}
	queue.Unlock()

	datum, err := spiderItemQueue.Get()
	if err != nil {
		return nil, err
	}

	ctx, ok := datum.(*data.Context)
	if !ok {
		return nil, nil
	}

	return ctx, nil
}

func (queue *Queues) Finish(ctx *data.Context) {
	nodeRequestQueue, ok := queue.NodeRequestMap[ctx.NodeName]
	if !ok {
		return
	}
	delete(nodeRequestQueue, ctx.Unique())
}
