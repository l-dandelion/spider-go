package sourceparser

import "errors"

var (
	ErrGetParsersSource = errors.New(`Can't get source from model.Rule["GenParsers"]`)
	ErrConvertToFunc = errors.New("Can't convert f to func()[]modeule.ParseResponse")
)