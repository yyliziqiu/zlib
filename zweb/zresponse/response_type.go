package zresponse

import (
	"net/http"

	"github.com/yyliziqiu/zlib/zerror"
)

/**
错误码定义
A开头 客户端错误
B开头 服务端错误
C开头 三方服务错误

0-99      HTTP 协议定义的错误
100-999   框架定义的错误
1000-9999 用户自定义错误
*/

var (
	BadRequestError          = zerror.New("A0001", "Bad Request").StatusCode(http.StatusBadRequest)
	UnauthorizedError        = zerror.New("A0002", "Unauthorized").StatusCode(http.StatusUnauthorized)
	ForbiddenError           = zerror.New("A0003", "Forbidden").StatusCode(http.StatusForbidden)
	NotFoundError            = zerror.New("A0004", "Not Found").StatusCode(http.StatusNotFound)
	MethodNotAllowedError    = zerror.New("A0005", "Method Not Allowed").StatusCode(http.StatusMethodNotAllowed)
	InternalServerErrorError = zerror.New("B0001", "Internal Server Error").StatusCode(http.StatusInternalServerError)
)

type ErrorResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResult(code string, message string) ErrorResult {
	return ErrorResult{
		Code:    code,
		Message: message,
	}
}

func NewErrorResultWithError(err *zerror.Error) ErrorResult {
	return NewErrorResult(err.Code, err.Message)
}

func errorResponse(err error, verbose bool) (int, ErrorResult) {
	var (
		statusCode = http.StatusBadRequest
		code       = BadRequestError.Code
		message    = BadRequestError.Message
	)

	zerr, ok := err.(*zerror.Error)
	if ok {
		statusCode, code, message = zerr.HTTP()
	} else if verbose {
		message = err.Error()
	}

	return statusCode, NewErrorResult(code, message)
}
