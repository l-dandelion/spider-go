package processors

import (
	"github.com/l-dandelion/yi-ants-go/core/module"
	"github.com/l-dandelion/yi-ants-go/lib/constant"
	"github.com/l-dandelion/yi-ants-go/core/processors/mysqlprocessor"
	"github.com/l-dandelion/yi-ants-go/core/processors/sourceprocessor"
	"github.com/l-dandelion/yi-ants-go/core/processors/model"
	"github.com/l-dandelion/yi-ants-go/core/processors/consoleprocessor"
)


func GenProcessorsByModel(model *model.Model) ([]module.ProcessItem, *constant.YiError){
	switch model.Type {
	case "mysql":
		return []module.ProcessItem{mysqlprocessor.DefaultMysqlProcessor}, nil
	case "console":
		return []module.ProcessItem{consoleprocessor.DefaultConsoleProcessor}, nil
	case "source":
		return sourceprocessor.GetSourceProcessorsFromModel(model)
	default:
		return nil, constant.NewYiErrorf(constant.ERR_UNSUPPORTED_MODEL_TYPE, "Unsupported model type.(modelType: %s)", model.Type)
	}
}

func GenProcessorsByModels(models []*model.Model) ([]module.ProcessItem, *constant.YiError){
	processors := []module.ProcessItem{}
	for _, model := range models {
		ps, yierr := GenProcessorsByModel(model)
		if yierr != nil {
			return nil, yierr
		}
		processors = append(processors, ps...)
	}
	return processors, nil
}