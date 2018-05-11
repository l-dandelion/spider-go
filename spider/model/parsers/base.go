package parsers

import (
	"github.com/l-dandelion/spider-go/spider/module"
	"github.com/l-dandelion/spider-go/spider/model/parsers/templateparser"
	"github.com/l-dandelion/spider-go/spider/model/parsers/sourceparser"
	"github.com/l-dandelion/spider-go/spider/model/parsers/model"
	"fmt"
)

func GenParsersByModel(model *model.Model) ([]module.ParseResponse, error){
	switch model.Type {
	case "template":
		return []module.ParseResponse{templateparser.GenTemplateParser(model)}, nil
	case "source":
		return sourceparser.GetSourceParsersFromModel(model)
	default:
		return nil, fmt.Errorf("Unsupported model type.(modelType: %s)", model.Type)
	}
}

func GenParsersByModels(models []*model.Model) ([]module.ParseResponse, error){
	parsers := []module.ParseResponse{}
	for _, model := range models {
		ps, err := GenParsersByModel(model)
		if err != nil {
			return nil, err
		}
		parsers = append(parsers, ps...)
	}
	return parsers, nil
}