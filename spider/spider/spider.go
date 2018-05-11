package spider

import "github.com/l-dandelion/spider-go/spider/scheduler"

type Spider struct {
	Job
	scheduler.Scheduler
}

func NewSpider(job *Job) *Spider {
	return &Spider{
		Job: *job,
	}
}