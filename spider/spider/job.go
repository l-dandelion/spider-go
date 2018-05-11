package spider

import (
	parsermodel "github.com/l-dandelion/spider-go/spider/model/parsers/model"
	processormodel "github.com/l-dandelion/spider-go/spider/model/processors/model"
	"time"
)

type Job struct {
	Name            string
	MaxDepth        uint32
	AcceptedDomains []string
	ParserModels    []*parsermodel.Model
	ProcessorModels []*processormodel.Model
	InitialUrls     []string
	CreatedAt       time.Time
}

func NewJob(name string,
	maxDepth uint32,
	AcceptedDomains []string,
	parserModels []*parsermodel.Model,
	processorModels []*processormodel.Model,
	initialUrls []string) *Job {
	return &Job{
		Name:            name,
		MaxDepth:        maxDepth,
		AcceptedDomains: AcceptedDomains,
		ParserModels:    parserModels,
		ProcessorModels: processorModels,
		InitialUrls:     initialUrls,
		CreatedAt:       time.Now(),
	}
}
