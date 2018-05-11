package sourceprocessor

import (
	"github.com/l-dandelion/yi-ants-go/core/module"
	"github.com/l-dandelion/yi-ants-go/lib/constant"
	"github.com/l-dandelion/yi-ants-go/lib/library/plugin"
	"github.com/l-dandelion/yi-ants-go/core/processors/model"
	"fmt"
)

func GetSourceProcessorsFromModel(model *model.Model) ([]module.ProcessItem, *constant.YiError) {
	fmt.Println(model)
	source, ok := model.Rule["GenProcessors"]
	if !ok {
		return nil, constant.NewYiErrorf(constant.ERR_GET_PROCESSORS_SOURCE, `Can't get source from model.Rule["GenProcessors"]`)
	}
	f, err := plugin.GenFuncFromStr(source, "GenProcessors")
	if err != nil {
		return nil, constant.NewYiErrore(constant.ERR_GET_PROCESSORS, err)
	}
	genFunc, ok := f.(func() []module.ProcessItem)
	if !ok {
		return nil, constant.NewYiErrorf(constant.ERR_GET_PROCESSORS, "Can't convert f to func()[]modeule.ProcessItem")
	}
	result := genFunc()
	fmt.Println(result)
	return result, nil
}
