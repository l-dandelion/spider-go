package consoleprocessor

import (
	"github.com/l-dandelion/yi-ants-go/core/module/data"
	"github.com/l-dandelion/yi-ants-go/lib/constant"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

var DefaultConsoleProcessor = func(item data.Item) (result data.Item, yierr *constant.YiError) {
	bytes, _ := json.Marshal(item)
	log.Info("Pipline :", string(bytes))
	return nil, nil
}