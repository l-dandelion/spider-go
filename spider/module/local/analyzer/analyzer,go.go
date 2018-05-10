package analyzer

import (
	"github.com/l-dandelion/spider-go/spider/module"
	"github.com/l-dandelion/spider-go/spider/module/stub"
	"github.com/l-dandelion/spider-go/spider/module/data"
	log "github.com/sirupsen/logrus"
	"github.com/l-dandelion/yi-ants-go/lib/library/reader"
	"net/url"
)

/*
 * create an instance for analyzer
 */
func New() (analyzer module.Analyzer) {
	moduleBase := stub.NewModuleInternal()
	return &myAnalyzer{
		ModuleInternal: moduleBase,
	}
}


/*
 * implementation of interface module.Analyzer
 */
type myAnalyzer struct {
	stub.ModuleInternal                        // module internal instance
}

/*
 * analyze the response and return the data list, or error list
 */
func (analyzer *myAnalyzer) Analyze(ctx *data.Context, parsers []module.ParseResponse) {
	analyzer.IncrHandlingNumber()
	defer analyzer.DecrHandlingNumber()
	analyzer.IncrCalledCount()

	// check args
	if len(parsers) == 0 {
		ctx.PushError(ErrEmptyParsers)
		return
	}
	if ctx.HTTPResp() == nil {
		ctx.PushError(ErrNilHTTPResponse)
		return
	}
	analyzer.IncrAcceptedCount()

	httpReq := ctx.HttpReq
	var reqURL *url.URL
	if httpReq != nil {
		reqURL = httpReq.URL
	}
	respDepth := ctx.Depth
	log.Infof("Parse the response (URL: %s, depth: %d)... \n", reqURL, respDepth)

	httpResp := ctx.HTTPResp()
	if httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	multiReader, err := reader.NewMultipleReader(httpResp.Body)
	if err != nil {
		ctx.PushError(err)
		return
	}
	for _, respParser := range parsers {
		if httpResp.Body != nil {
			httpResp.Body.Close()
		}
		httpResp.Body = multiReader.Reader()
		respParser(ctx)
	}
	if len(ctx.ErrorList) == 0 {
		analyzer.IncrCompletedCount()
	}
}