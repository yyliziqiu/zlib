package zerror

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    string
	Message string
}

func New(code string, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

func (e *Error) With(v interface{}) *Error {
	err := &Error{Code: e.Code}
	switch v.(type) {
	case error:
		err.Message = v.(error).Error()
	case string:
		err.Message = v.(string)
	default:
		err.Message = fmt.Sprintf("%v", v)
	}
	return err
}

func (e *Error) Wrap(err error) *Error {
	return &Error{
		Code:    e.Code,
		Message: fmt.Sprintf("%s [%v]", e.Message, err),
	}
}

func (e *Error) Format(message string, a ...interface{}) *Error {
	return &Error{
		Code:    e.Code,
		Message: fmt.Sprintf(message, a...),
	}
}

func (e *Error) Fields(a ...interface{}) *Error {
	return &Error{
		Code:    e.Code,
		Message: fmt.Sprintf(e.Message, a...),
	}
}

func (e *Error) HTTP() (int, string, string) {
	statusCode := http.StatusBadRequest
	if e.Code[0] != 'A' {
		statusCode = http.StatusInternalServerError
	}
	return statusCode, e.Code, e.Message
}
