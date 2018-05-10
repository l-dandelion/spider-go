package module

import (
	"github.com/l-dandelion/spider-go/spider/module/data"
)

/*
 * the type of internal counts of module.
 * include: called count, accepted count, completed count, handling number
 */
type Counts struct {
	CalledCount    uint64 // 调用计数
	AcceptedCount  uint64 // 接受计数
	CompletedCount uint64 // 完成计数
	HandlingNumber uint64 // 正在处理的数目
}

/*
 * the struct of summary of module
 */
type SummaryStruct struct {
	Called    uint64      `json:"called"`          // 调用计数
	Accepted  uint64      `json:"accepted"`        // 接受计数
	Completed uint64      `json:"completed"`       // 完成计数
	Handling  uint64      `json:"handling"`        // 正在处理的数目
	Extra     interface{} `json:"extra,omitempty"` // 额外信息
}

/*
 * basic interface representing module
 * the implementation type of the interface must be concurrent and secure.
 */
type Module interface {
	CalledCount() uint64    // 获取已调用的数目
	AcceptedCount() uint64  // 获取已接受的数目
	CompletedCount() uint64 // 获取已完成的数目
	HandlingNumber() uint64 // 获取正在处理的数目
	Counts() Counts         // 获取各个计数的结构体
	Summary() SummaryStruct // 获取摘要
}

/*
 * interface for the downloader
 * inherit from module
 * the implementation type of the interface must be concurrent and secure.
 */
type Downloader interface {
	Module                      // inherit from module
	Download(ctx *data.Context) // 根据请求下载相应页面，将响应存入ctx中
}

/*
 * function for parese response
 */
type ParseResponse func(ctx *data.Context)

/*
 * interface for the analyzer
 * inherit from module
 * the implementation type of the interface must be concurrent and secure.
 */
type Analyzer interface {
	Analyze(ctx *data.Context, parsers []ParseResponse) // 根据解析规则解析响应，并将解析结果存入ctx中
}

/*
 * func for process the item
 */
type ProcessItem func(item data.Item) (result data.Item, err error)

/*
 * interface for the pipeline
 * inherit from module
 * the implementation type of the interface must be concurrent and secure.
 */
type Pipeline interface {
	Module
	Send(ctx *data.Context, processors []ProcessItem) // Send item to the pipeline, and process by item processors one by one
}