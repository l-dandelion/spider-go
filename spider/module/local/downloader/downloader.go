package downloader

import (
	"net/http"

	"github.com/l-dandelion/spider-go/spider/module"
	"github.com/l-dandelion/spider-go/spider/module/data"
	"github.com/l-dandelion/spider-go/spider/module/stub"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

var (
	RetryTimes = 3
	DefaultClient = &http.Client{
		Transport: &http.Transport{
		DialContext: (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 10 * time.Second,
		DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   5,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		},
	}
	DefaultDownloader = New(nil)
)

/*
 * create an instance for downloader
 */
func New(client *http.Client) (downloader module.Downloader) {
	moduleBase := stub.NewModuleInternal()
	if client == nil {
		client = DefaultClient
	}

	return &myDownloader{
		ModuleInternal: moduleBase,
		httpClient:     client,
	}
}

/*
 * implementation of interface module.Downloader
 */
type myDownloader struct {
	stub.ModuleInternal              //module internal instance
	httpClient          *http.Client //http client for downloading
}

/*
 * download according to request, return a response if success, or an error return
 */
func (downloader *myDownloader) Download(ctx *data.Context) {
	downloader.IncrHandlingNumber()
	defer downloader.DecrHandlingNumber()
	downloader.IncrCalledCount()

	if ctx.HttpReq == nil {
		ctx.PushError(ErrNilHTTPRequest)
		return
	}

	downloader.IncrAcceptedCount()

	var (
		httpResp *http.Response
		err      error
	)

	log.Infof("Do the request (URL: %s, depth: %d)... \n", ctx.HttpReq.URL, ctx.Depth)

	// try to download RETRY_TIMES times
	for i := 0; i < RetryTimes; i++ {
		httpResp, err = downloader.httpClient.Do(ctx.HttpReq)
		if err == nil {
			break
		}
	}

	if err != nil {
		ctx.PushError(err)
		return
	}
	ctx.SetResponse(httpResp)
	downloader.IncrCompletedCount()
}
