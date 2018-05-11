package sourceparser

import (
	"github.com/l-dandelion/spider-go/spider/module"
	"github.com/l-dandelion/spider-go/lib/library/plugin"
	"github.com/l-dandelion/spider-go/spider/model/parsers/model"
)

func GetSourceParsersFromModel(model *model.Model) ([]module.ParseResponse, error) {
	source, ok := model.Rule["GenParsers"]
	if !ok {
		return nil, ErrGetParsersSource
	}
	f, err := plugin.GenFuncFromStr(source, "GenParsers")
	if err != nil {
		return nil, err
	}
	genFunc, ok := f.(func() []module.ParseResponse)
	if !ok {
		return nil, ErrConvertToFunc
	}
	result := genFunc()
	return result, nil
}
