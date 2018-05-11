package processors

import (
	"github.com/l-dandelion/spider-go/spider/module"
	"github.com/l-dandelion/spider-go/spider/model/processors/mysqlprocessor"
	"github.com/l-dandelion/spider-go/spider/model/processors/sourceprocessor"
	"github.com/l-dandelion/spider-go/spider/model/processors/model"
	"github.com/l-dandelion/spider-go/spider/model/processors/consoleprocessor"
	"fmt"
)


func GenProcessorsByModel(model *model.Model) ([]module.ProcessItem, error){
	switch model.Type {
	case "mysql":
		return []module.ProcessItem{mysqlprocessor.DefaultMysqlProcessor}, nil
	case "console":
		return []module.ProcessItem{consoleprocessor.DefaultConsoleProcessor}, nil
	case "source":
		return sourceprocessor.GetSourceProcessorsFromModel(model)
	default:
		return nil, fmt.Errorf("Unsupported model type.(modelType: %s)", model.Type)
	}
}

func GenProcessorsByModels(models []*model.Model) ([]module.ProcessItem, error){
	processors := []module.ProcessItem{}
	for _, model := range models {
		ps, err := GenProcessorsByModel(model)
		if err != nil {
			return nil, err
		}
		processors = append(processors, ps...)
	}
	return processors, nil
}