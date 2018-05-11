package templateparser

import (
	"github.com/l-dandelion/yi-ants-go/core/module"
	"github.com/l-dandelion/yi-ants-go/core/module/data"
	"github.com/l-dandelion/yi-ants-go/core/parsers/filter"
	"github.com/l-dandelion/yi-ants-go/lib/constant"
	"github.com/l-dandelion/yi-ants-go/core/parsers/model"
)

func GenTemplateParser(model *model.Model) module.ParseResponse {
	return func(resp *data.Response) ([]data.Data, []*constant.YiError) {
		if len(model.AcceptedRegUrls) > 0 && !filter.Filter(resp.HTTPRequest().URL.String(), model.AcceptedRegUrls) {
			return nil, nil
		}
		return TemplateRuleProcess(model, resp)
	}
}