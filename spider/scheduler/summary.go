package scheduler

import (
	"github.com/l-dandelion/spider-go/spider/module"
	"github.com/l-dandelion/yi-ants-go/lib/constant"
)

type Summary struct {
	MaxDepth       uint32
	DownloadCounts module.SummaryStruct
	Status         int8
	StatusDesc     string
	NumURL         uint64
}

func (sched *myScheduler) Summary() Summary {
	status := sched.Status()
	return Summary{
		MaxDepth:       sched.maxDepth,
		DownloadCounts: sched.downloaderCounts.Summary(),
		Status:         status,
		StatusDesc:     constant.GlobalArrRunningStatusDesc[status],
		NumURL:         sched.urlMap.Len(),
	}
}