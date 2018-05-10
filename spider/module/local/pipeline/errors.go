package pipeline

import (
	"errors"
)

var (
	ErrEmptyProcessors = errors.New("Empty processor list")
)
