package mysqlprocessor

import (
	"testing"
	"fmt"
)

func TestGensql(t *testing.T) {
	model1 := NewDBModel("test", map[string]interface{}{
		"tags": "123",
		"href": "baidu.com",
	})
	model2 := NewDBModel("test", map[string]interface{}{
		"tags": "1234",
		"href": "baidu4.com",
	})
	models := []*DBModel{model1, model2}
	fmt.Println(GenInsertModelsSql(models))
	fmt.Println(GenInsertModelsArgs(models))
}
