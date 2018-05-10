package plugin

import (
	"github.com/l-dandelion/yi-ants-go/lib/utils"
	"os/exec"
	"errors"
	"plugin"
	"math/rand"
	"fmt"
)

var path = "../temp"

func GenFuncFromStr(sourceStr, funcName string) (interface{}, error) {
	fileName := fmt.Sprintf("%x", rand.Uint64()) + ".go"
	dirPath := path
	err := utils.SaveFile(dirPath, fileName, []byte(sourceStr))
	if err != nil {
		return nil, err
	}
	return GenFuncFromSource(dirPath+"/"+fileName, funcName)
}

func GenFuncFromSource(filePath, funcName string) (interface{}, error) {
	fileName := fmt.Sprintf("%x", rand.Uint64()) + ".so"
	dirPath := path
	cmd := exec.Command("go", "build", "-o", dirPath + "/" + fileName, "-buildmode=plugin", filePath)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if string(out) != "" {
		return nil, errors.New(string(out))
	}
	p, err := plugin.Open(dirPath+"/"+fileName)
	if err != nil {
		return nil, err
	}
	f, err := p.Lookup(funcName)
	if err != nil {
		return nil, err
	}
	return f, err
}