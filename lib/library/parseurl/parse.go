package parseurl

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

func ParseReqUrl(urls []string, ctx map[string]interface{}) []string {
	resultUrls := []string{}
	for _, url := range urls {
		results, ok := isRuleReq(url, ctx)
		if ok {
			resultUrls = append(resultUrls, results...)
		} else {
			resultUrls = append(resultUrls, url)
		}
	}
	return resultUrls
}

func FindRule(text string) [][]string {
	reg := regexp.MustCompile(`{([^}]+)}`)
	return reg.FindAllStringSubmatch(text, -1)
}

func isRuleReq(reqUrl string, ctx map[string]interface{}) ([]string, bool) {
	reqUrls := []string{reqUrl}
	outUrls := []string{}
	finalUrls := []string{}
	isMatch := false

	rules := FindRule(reqUrl)
	if len(rules) > 0 {
		isMatch = true
	} else {
		return nil, false
	}

	if ctx != nil {
		reqUrls, isMatch = PraseParamCtx(reqUrl, ctx)
	}
	for _, r := range reqUrls {
		outUrls = append(outUrls, PraseOffset(r)...)
	}

	for _, r := range outUrls {
		finalUrls = append(finalUrls, PraseOr(r)...)
	}

	if isMatch {
		return finalUrls, true
	}

	return finalUrls, isMatch
}

// http://xxxxxxxx.com/abc/{begin-end,offset}/   example:{1-400,10}
func PraseOffset(reqUrl string) []string {
	reqUrls := []string{}
	outUrls := []string{}

	rules := FindRule(reqUrl)
	if len(rules) <= 0 {
		return []string{reqUrl}
	}

	var begin, end, offset int
	var rule string
	for _,rulee :=range rules{
		rule = rulee[1]
		sp := strings.Split(rule, ",")

		if len(sp) != 2 {
			continue
		}

		rs := strings.Split(sp[0], "-")

		var err error
		begin, err = strconv.Atoi(rs[0])
		end, err = strconv.Atoi(rs[1])
		offset, err = strconv.Atoi(sp[1])
		if err != nil {
			continue
		}
		if offset == 0 {
			continue
		}

		break
	}

	if begin == 0 && end == 0 && offset == 0{
		return []string{reqUrl}
	}

	for i := begin; i <= end; i = i + offset {
		url := strings.Replace(reqUrl, "{"+rule+"}", strconv.Itoa(i), 1)
		reqUrls = append(reqUrls, url)
	}

	for _, r := range reqUrls {
		outUrls = append(outUrls, PraseOffset(r)...)
	}

	return outUrls
}

// http://xxxxxxxx.com/abc/{id1|id2|id3}/
func PraseOr(reqUrl string) []string {
	reqUrls := []string{}
	outUrls := []string{}

	rules := FindRule(reqUrl)
	if len(rules) <= 0 {
		return []string{reqUrl}
	}
	ruleArray := rules[0]
	rule := ruleArray[1]
	sp := strings.Split(rule, "|")
	if len(sp) < 2 {
		return []string{reqUrl}
	}

	for _, word := range sp {
		url := strings.Replace(reqUrl, "{"+rule+"}", word, 1)
		reqUrls = append(reqUrls, url)
	}

	for _, r := range reqUrls {
		outUrls = append(outUrls, PraseOr(r)...)
	}

	return outUrls
}

// http://xxxxxxxx.com/abc/{name}/{id}/
func PraseParamCtx(reqUrl string, ctx map[string]interface{}) ([]string, bool) {
	urls := []string{}

	count := strings.Count(reqUrl, "$")
	if count <= 0 {
		return urls, false
	}

	for ctxName, ruleUrl := range ctx {
		urlStr, ok := ruleUrl.(string)

		if ok {
			reqUrl = strings.Replace(reqUrl, "{$"+ctxName+"}", string(urlStr), -1)
			//reqUrl = strings.Replace(reqUrl, "$"+ctxName, string(urlStr), -1)
		}

		urlNumber, ok := ruleUrl.(json.Number)

		if ok {
			reqUrl = strings.Replace(reqUrl, "{$"+ctxName+"}", string(urlNumber), -1)
			//reqUrl = strings.Replace(reqUrl, "$"+ctxName, string(urlNumber), -1)
		}

		urlInt, ok := ruleUrl.(int)
		if ok {
			reqUrl = strings.Replace(reqUrl, "{$"+ctxName+"}", strconv.Itoa(urlInt), -1)
			//reqUrl = strings.Replace(reqUrl, "$"+ctxName, strconv.Itoa(urlInt), -1)
		}
	}
	urls = append(urls, reqUrl)
	return urls, true
}
