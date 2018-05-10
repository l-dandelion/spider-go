package pipeline

import (
	"github.com/l-dandelion/spider-go/spider/module"
	"github.com/l-dandelion/spider-go/spider/module/data"
	"github.com/l-dandelion/spider-go/spider/module/stub"
	log "github.com/sirupsen/logrus"
)

/*
 * implementation of interface module.Pipeline
 */
type myPipeline struct {
	stub.ModuleInternal
}

/*
 * create an instance for module.Pipeline
 */
func New() (pipeline module.Pipeline) {
	moduleBase := stub.NewModuleInternal()
	return &myPipeline{
		ModuleInternal: moduleBase,
	}
}

/*
 * process item
 */
func (pipeline *myPipeline) Send(ctx *data.Context, processors []module.ProcessItem) {
	pipeline.IncrHandlingNumber()
	defer pipeline.DecrHandlingNumber()
	pipeline.IncrCalledCount()
	if len(processors) == 0 {
		ctx.PushError(ErrEmptyProcessors)
		return
	}
	pipeline.IncrAcceptedCount()
	if len(ctx.ItemList) > 0 {
		for _, item := range ctx.ItemList {
			log.Infof("Process item %+v... \n", item)
			currentItem := item
			for _, processor := range processors {
				processedItem, err := processor(currentItem)
				if err != nil {
					ctx.PushError(err)
				}
				if processedItem != nil {
					currentItem = processedItem
				}
			}
		}
	}
	if len(ctx.ErrorList) == 0 {
		pipeline.IncrCompletedCount()
	}
	return
}