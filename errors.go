package sfox

import (
	"errors"
	"fmt"
)

var (
	ErrUnknownPayload = errors.New("Unknown payload")
)

type ErrHttp struct {
	StatusCode int
	Text       string
}

func (e ErrHttp) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Text)
}

type ErrHttpClient struct {
	ErrHttp
}

type ErrHttpServer struct {
	ErrHttp
}
