package zresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/zlib/zerror"
)

func Response(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, data)
}

func ResponseError(ctx *gin.Context, statusCode int, code string, message string) {
	ctx.JSON(statusCode, NewErrorResult(code, message))
}

func OK(ctx *gin.Context) {
	ctx.String(http.StatusOK, "")
}

func Result(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func Error(ctx *gin.Context, err error) {
	ctx.JSON(errorResponse(err))
}

func errorResponse(err error) (int, ErrorResult) {
	var (
		statusCode = http.StatusBadRequest
		code       = BadRequestError.Code
		message    = err.Error()
	)

	zerr, ok := err.(*zerror.Error)
	if ok {
		statusCode, code, message = zerr.HTTP()
	}

	return statusCode, NewErrorResult(code, message)
}

func ErrorString(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, NewErrorResult(BadRequestError.Code, message))
}

func AbortOK(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}

func AbortResult(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, data)
}

func AbortError(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(errorResponse(err))
}

func AbortErrorString(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResult(BadRequestError.Code, message))
}

func AbortBadRequest(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResultWithError(BadRequestError))
}

func AbortUnauthorized(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, NewErrorResultWithError(UnauthorizedError))
}

func AbortForbidden(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusForbidden, NewErrorResultWithError(ForbiddenError))
}

func AbortNotFound(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusNotFound, NewErrorResultWithError(NotFoundError))
}

func AbortMethodNotAllowed(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, NewErrorResultWithError(MethodNotAllowedError))
}

func AbortInternalServerError(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResultWithError(InternalServerErrorError))
}
