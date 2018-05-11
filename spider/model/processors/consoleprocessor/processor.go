package consoleprocessor

import (
	"github.com/l-dandelion/spider-go/spider/module/data"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

var DefaultConsoleProcessor = func(item data.Item) (result data.Item, err error) {
	bytes, _ := json.Marshal(item)
	log.Info("Pipline :", string(bytes))
	return nil, nil
}