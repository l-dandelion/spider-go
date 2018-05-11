package templateparser

import (
	"github.com/l-dandelion/spider-go/spider/module"
	"github.com/l-dandelion/spider-go/spider/module/data"
	"github.com/l-dandelion/spider-go/spider/model/parsers/filter"
	"github.com/l-dandelion/spider-go/spider/model/parsers/model"
)

func GenTemplateParser(model *model.Model) module.ParseResponse {
	return func(ctx *data.Context) {
		if len(model.AcceptedRegUrls) > 0 && !filter.Filter(ctx.HttpReq.URL.String(), model.AcceptedRegUrls) {
			return
		}
		TemplateRuleProcess(model, ctx)
	}
}