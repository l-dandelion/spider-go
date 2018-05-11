package sourceparser

import (
	"github.com/l-dandelion/yi-ants-go/core/module"
	"github.com/l-dandelion/yi-ants-go/lib/constant"
	"github.com/l-dandelion/yi-ants-go/lib/library/plugin"
	"github.com/l-dandelion/yi-ants-go/core/parsers/model"
)

func GetSourceParsersFromModel(model *model.Model) ([]module.ParseResponse, *constant.YiError) {
	source, ok := model.Rule["GenParsers"]
	if !ok {
		return nil, constant.NewYiErrorf(constant.ERR_GET_PARSERS_SOURCE, `Can't get source from model.Rule["GenParsers"]`)
	}
	f, err := plugin.GenFuncFromStr(source, "GenParsers")
	if err != nil {
		return nil, constant.NewYiErrore(constant.ERR_GET_PARSERS, err)
	}
	genFunc, ok := f.(func() []module.ParseResponse)
	if !ok {
		return nil, constant.NewYiErrorf(constant.ERR_GET_PARSERS, "Can't convert f to func()[]modeule.ParseResponse")
	}
	result := genFunc()
	return result, nil
}
