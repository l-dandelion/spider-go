package sourceprocessor

import (
	"errors"
	"github.com/l-dandelion/spider-go/lib/library/plugin"
	"github.com/l-dandelion/spider-go/spider/model/processors/model"
	"github.com/l-dandelion/spider-go/spider/module"
)

func GetSourceProcessorsFromModel(model *model.Model) ([]module.ProcessItem, error) {
	source, ok := model.Rule["GenProcessors"]
	if !ok {
		return nil, errors.New(`Can't get source from model.Rule["GenProcessors"]`)
	}
	f, err := plugin.GenFuncFromStr(source, "GenProcessors")
	if err != nil {
		return nil, err
	}
	genFunc, ok := f.(func() []module.ProcessItem)
	if !ok {
		return nil, errors.New("Can't convert f to func()[]modeule.ProcessItem")
	}
	result := genFunc()
	return result, nil
}
