package data

import "net/http"

/*
 * request struct
 */
type Request struct {
	NodeName   string
	SpiderName string
	HttpReq    *http.Request          // the http request
	Depth      uint32                 // crawl depth
	Proxy      string                 // use proxy if not empty
	Extra      map[string]interface{} // additional information(used for context)
}

/*
 * set extra infomation
 */
func (req *Request) SetExtra(key string, val interface{}) *Request {
	req.Extra[key] = val
	return req
}

/*
 * check the request
 */
func (req *Request) Valid() bool {
	return req.HttpReq != nil && req.HttpReq.URL != nil
}

/*
 * New an instance of Request
 */
func NewRequest(httpReq *http.Request, extras ...map[string]interface{}) *Request {
	var extra map[string]interface{}
	if len(extras) != 0 {
		extra = extras[0]
	} else {
		extra = map[string]interface{}{}
	}
	return &Request{
		HttpReq: httpReq,
		Extra:   extra,
	}
}

/*
 * add cookie
 */
func (req *Request) AddCookie(key, value string) *Request {
	c := &http.Cookie{
		Name:  key,
		Value: value,
	}
	req.HttpReq.AddCookie(c)
	return req
}

/*
 * set header
 */
func (req *Request) SetHeader(key, value string) *Request {
	req.HttpReq.Header.Set(key, value)
	return req
}

/*
 * set user agent
 */
func (req *Request) SetUserAgent(ua string) *Request {
	req.SetHeader("User-Agent", ua)
	return req
}

/*
 * set referer
 */
func (req *Request) SetReferer(referer string) *Request {
	req.SetHeader("referer", referer)
	return req
}

/*
 * set proxy
 */
func (req *Request) SetProxy(proxy string) *Request {
	req.Proxy = proxy
	return req
}

/*
 * set spider name
 */
func (req *Request) SetSpiderName(name string) *Request {
	req.SpiderName = name
	return req
}

/*
 * set node name
 */
func (req *Request) SetNodeName(nodeName string) *Request {
	req.NodeName = nodeName
	return req
}

/*
 * set crawl depth
 */
func (req *Request) SetDepth(depth uint32) *Request {
	req.Depth = depth
	return req
}
