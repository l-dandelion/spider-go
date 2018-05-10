package analyzer

import (
	"errors"
)

var (
	ErrNilHTTPResponse = errors.New("Nil HTTP response.")
	ErrEmptyParsers = errors.New("Empty parser list")
)
