package mysqlprocessor

import (
	"github.com/l-dandelion/yi-ants-go/core/module/data"
	"github.com/l-dandelion/yi-ants-go/lib/constant"
)

var DefaultMysqlProcessor = func(item data.Item) (result data.Item, yierr *constant.YiError) {
	if _, ok := item["kind"]; !ok {
		return nil, constant.NewYiErrorf(constant.ERR_CRAWL_PIPELINE, "Kind is empty.")
	}
	kind := item["kind"].(string)
	if kind != "mysql" {
		return nil, nil
	}
	if _, ok := item["tableName"]; !ok {
		return nil, constant.NewYiErrorf(constant.ERR_CRAWL_PIPELINE, "Table name not found.")
	}
	tableName := item["tableName"].(string)
	delete(item, "tableName")
	delete(item, "kind")
	dbModel := NewDBModel(tableName, item)
	item["tableName"] = tableName
	item["kind"] = kind
	AddPrepare(dbModel)
	return nil, nil
}