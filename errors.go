package memgo

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrorNotStored = errors.New("memcached returned NOT_STORED")
	ErrCacheMiss   = errors.New("memcached cache miss")
)

type ErrClient struct {
	Response string
}

func (e *ErrClient) Error() string {
	return fmt.Sprintf("client error: %s", e.Response)
}

type ErrServer struct {
	Response string
}

func (e *ErrServer) Error() string {
	return fmt.Sprintf("server error: %s", e.Response)
}

func handleErrorResponse(response string) error {
	if strings.HasPrefix(response, "CLIENT_ERROR") {
		return &ErrClient{Response: response}
	} else {
		return &ErrServer{Response: response}
	}
}
