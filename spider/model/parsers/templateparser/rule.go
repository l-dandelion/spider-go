package templateparser

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"github.com/l-dandelion/spider-go/spider/module/data"
	"github.com/l-dandelion/yi-ants-go/lib/utils"
	"net/http"
	"github.com/l-dandelion/yi-ants-go/lib/library/parseurl"
	"github.com/l-dandelion/spider-go/spider/model/parsers/model"
)

func TemplateRuleProcess(model *model.Model, ctx *data.Context) {
	resp := ctx.Response
	rule := model.Rule

	doc, err := resp.GetDom()
	if err != nil {
		ctx.PushError(err)
		return
	}

	if len(model.WantedRegUrls) > 0 {
		doc.Find("a").Each(func(i int, sel *goquery.Selection) {
			href, _ := sel.Attr("href")
			href, err = utils.GetComplateUrl(ctx.HttpReq.URL, href)
			if err != nil {
				ctx.PushError(err)
				return
			}
			
			httpReq, err := http.NewRequest("GET", href, nil)
			if err != nil {
				ctx.PushError(err)
				return
			}

			ctx.PushRequest(data.NewRequest(httpReq))
		})
	}

	resultType := "map"
	rootSel := ""

	v, ok := rule["node"]
	if ok {
		contentInfo := strings.Split(v, "|")
		resultType = contentInfo[0]
		rootSel = contentInfo[1]
	}

	if resultType == "array" {
		doc.Find(rootSel).Each(func(i int, s *goquery.Selection) {
			mdata := getMapFromDom(rule, s)
			if mdata == nil {
				return
			}
			ctx.PushItem(data.Item(mdata))
			if len(model.AddQueue) > 0 {
				urls := parseurl.ParseReqUrl(model.AddQueue, mdata)
				for _, u := range urls {
					httpReq, err := http.NewRequest("GET", u, nil)
					if err != nil {
						ctx.PushError(err)
						return
					}
					ctx.PushRequest(data.NewRequest(httpReq))
				}
			}
		})
	}

	if resultType == "map" {
		mdata := getMapFromDom(rule, doc.Selection)
		ctx.PushItem(data.Item(mdata))
		if len(model.AddQueue) > 0 {
			urls := parseurl.ParseReqUrl(model.AddQueue, mdata)
			for _, u := range urls {
				httpReq, err := http.NewRequest("GET", u, nil)
				if err != nil {
					ctx.PushError(err)
					return
				}
				ctx.PushRequest(data.NewRequest(httpReq))
			}
		}
	}

	return
}

func getMapFromDom(rule map[string]string, node *goquery.Selection) map[string]interface{} {

	result := make(map[string]interface{})

	isNull := true

	for key, value := range rule {

		if key == "node" {
			continue
		}

		rules := strings.Split(value, "|")
		ValueType := strings.Split(rules[0], ".")

		if len(rules) < 2 {
			result[key] = value
			continue
		}

		s := node.Find(rules[1])
		switch ValueType[0] {
		case "text":
			result[key] = s.Text()
		case "html":
			result[key], _ = s.Html()
		case "attr":
			if len(ValueType) < 2 {
				continue
			}
			result[key], _ = s.Attr(ValueType[1])
		case "texts":
			arr := []string{}
			s.Each(func(i int, sel *goquery.Selection) {
				text := sel.Text()
				arr = append(arr, text)
			})
			j, _ := json.Marshal(arr)
			result[key] = string(j)
		case "htmls":
			arr := []string{}
			s.Each(func(i int, sel *goquery.Selection) {
				html, _ := s.Html()
				arr = append(arr, html)
			})
			j, _ := json.Marshal(arr)
			result[key] = string(j)
		case "attrs":
			arr := []string{}
			attr := ""
			s.Each(func(i int, sel *goquery.Selection) {
				if len(ValueType) >= 2 {
					attr, _ = sel.Attr(ValueType[1])
					arr = append(arr, attr)
				}
			})
			result[key] = arr
		default:
			result[key] = value
		}
		res, ok := result[key].(string)
		if ok || len(res) != 0 {
			isNull = false
		}
	}

	if isNull == true {
		return nil
	}

	return result
}