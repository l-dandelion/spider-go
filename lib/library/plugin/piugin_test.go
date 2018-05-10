package plugin

import "testing"

var funcStr = `package main
import "fmt"
func Hello() {
	fmt.Println("Hello")
}`

func TestPlugin(t *testing.T) {
	f, err := GenFuncFromStr(funcStr, "Hello")
	if err != nil {
		t.Fatal(err)
	}
	f.(func())()
}
