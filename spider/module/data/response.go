package data

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/l-dandelion/yi-ants-go/lib/library/reader"
	"io"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"mime"
	"strings"
	"bytes"
)

/*
 * Response
 * httpResp: the http response for thr request
 * req: the request for the response
 */
type Response struct {
	httpResp *http.Response    // the request for the response
	text     []byte            // body's []byte type
	dom      *goquery.Document // body's Dom type if body is html
}

/*
 * New an instance of Response
 */
func NewResponse(httpResp *http.Response) *Response {
	return &Response{httpResp: httpResp}
}

/*
 * get http response
 */
func (resp *Response) HTTPResp() *http.Response {
	return resp.httpResp
}

/*
 * check the response
 */
func (resp *Response) Valid() bool {
	return resp.httpResp != nil && resp.httpResp.Body != nil
}

/*
 * get body's []byte
 */
func (resp *Response) GetText() ([]byte, error) {
	if resp.text != nil {
		return resp.text, nil
	}
	multiReader, err := reader.NewMultipleReader(resp.httpResp.Body)
	resp.httpResp.Body = multiReader.Reader()
	defer func() {
		resp.httpResp.Body.Close()
		resp.httpResp.Body = multiReader.Reader()
	}()
	var contentType, pageEncode string

	// read firstly content-type from response header
	contentType = resp.httpResp.Header.Get("Content-Type")
	if _, params, err := mime.ParseMediaType(contentType); err == nil {
		if cs, ok := params["charset"]; ok {
			pageEncode = strings.ToLower(strings.TrimSpace(cs))
		}
	}

	// read content-type from request header
	if len(pageEncode) == 0 {
		contentType = resp.httpResp.Request.Header.Get("Content-Type")
		if _, params, err := mime.ParseMediaType(contentType); err == nil {
			if cs, ok := params["charset"]; ok {
				pageEncode = strings.ToLower(strings.TrimSpace(cs))
			}
		}
	}

	switch pageEncode {
	case "utf8", "utf-8", "unicode-1-1-utf-8":
	default:
		// get converter to utf-8
		// Charset auto determine. Use golang.org/x/net/html/charset. Get response body and change it to utf-8
		var destReader io.Reader

		if len(pageEncode) == 0 {
			destReader, err = charset.NewReader(resp.httpResp.Body, "")
		} else {
			destReader, err = charset.NewReaderLabel(pageEncode, resp.httpResp.Body)
		}

		if err == nil {
			resp.text, err = ioutil.ReadAll(destReader)
			if err == nil {
				return resp.text, err
			}
		}

	}
	resp.text, err = ioutil.ReadAll(resp.httpResp.Body)

	return resp.text, err
}


/*
 * get body's dom
 */
func (resp *Response) GetDom() (*goquery.Document, error) {
	if resp.dom != nil {
		return resp.dom, nil
	}
	text, err := resp.GetText()
	if err != nil {
		return nil, err
	}
	resp.dom, err = goquery.NewDocumentFromReader(bytes.NewReader(text))
	return resp.dom, err
}