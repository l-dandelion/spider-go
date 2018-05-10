package downloader

import (
	"errors"
)

var (
	ErrNilHTTPRequest = errors.New("Nil HTTP request.")
)