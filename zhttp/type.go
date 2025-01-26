package zhttp

import (
	"fmt"
	"net/http"
)

type JsonResponse interface {
	Failed() bool
}

type ResponseError struct {
	status int
	errstr string
}

func newResponseError(status int, errstr string) *ResponseError {
	return &ResponseError{
		status: status,
		errstr: errstr,
	}
}

func (e ResponseError) Status() int {
	return e.status
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("[%d] %s", e.status, e.errstr)
}

type HTTPLog struct {
	Method       string
	Request      *http.Request
	RequestBody  []byte
	Response     *http.Response
	ResponseBody []byte
	Error        error
	Cost         string
}
