package parsers

import (
	"github.com/l-dandelion/yi-ants-go/core/module"
	"github.com/l-dandelion/yi-ants-go/core/parsers/templateparser"
	"github.com/l-dandelion/yi-ants-go/lib/constant"
	"github.com/l-dandelion/yi-ants-go/core/parsers/sourceparser"
	"github.com/l-dandelion/yi-ants-go/core/parsers/model"
)

func GenParsersByModel(model *model.Model) ([]module.ParseResponse, *constant.YiError){
	switch model.Type {
	case "template":
		return []module.ParseResponse{templateparser.GenTemplateParser(model)}, nil
	case "source":
		return sourceparser.GetSourceParsersFromModel(model)
	default:
		return nil, constant.NewYiErrorf(constant.ERR_UNSUPPORTED_MODEL_TYPE, "Unsupported model type.(modelType: %s)", model.Type)
	}
}

func GenParsersByModels(models []*model.Model) ([]module.ParseResponse, *constant.YiError){
	parsers := []module.ParseResponse{}
	for _, model := range models {
		ps, yierr := GenParsersByModel(model)
		if yierr != nil {
			return nil, yierr
		}
		parsers = append(parsers, ps...)
	}
	return parsers, nil
}