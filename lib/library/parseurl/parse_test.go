package parseurl

import (
	"testing"
	"fmt"
)

func TestParse(t *testing.T) {

	urlStrs := []string{"{$img}"}
	mdata := map[string]interface{}{
		"img": "test",
	}
	fmt.Println(ParseReqUrl(urlStrs, mdata))
}