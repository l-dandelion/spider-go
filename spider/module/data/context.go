package data

import "net/http"

/*
 * Context 上下文，记录着请求、请求对应的响应、处理产生的条目以及错误信息
 */
type Context struct {
	Request
	Response
	ErrorList   []error
	RequestList []*Request
	ItemList    []Item
}

/*
 * 新建一个context
 */
func NewContext(req *Request) *Context {
	return &Context{
		Request: *req,
	}
}

/*
 * 设置响应
 */
func (ctx *Context) SetResponse(resp *http.Response) *Context {
	ctx.Response = *NewResponse(resp)
	return ctx
}

/*
 * 添加错误
 */
func (ctx *Context) PushError(errs ...error) *Context {
	if len(errs) == 0 {
		return ctx
	}
	if ctx.ErrorList == nil {
		ctx.ErrorList = []error{}
	}
	for _, err := range errs {
		ctx.ErrorList = append(ctx.ErrorList, err)
	}
	return ctx
}

/*
 * 添加新请求
 */
func (ctx *Context) PushRequest(reqs []*Request) *Context {
	if len(reqs) == 0 {
		return ctx
	}
	if ctx.ErrorList == nil {
		ctx.RequestList = []*Request{}
	}
	for _, req := range reqs {
		ctx.RequestList = append(ctx.RequestList, req)
	}
	return ctx

}

/*
 * 添加条目
 */
func (ctx *Context) PushItem(items []Item) *Context {
	if len(items) == 0 {
		return ctx
	}
	if ctx.ErrorList == nil {
		ctx.ItemList = []Item{}
	}
	for _, item := range items {
		ctx.ItemList = append(ctx.ItemList, item)
	}
	return ctx
}