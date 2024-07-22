package zweb

import (
	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/zlib/zerror"
	"github.com/yyliziqiu/zlib/zweb/zresponse"
)

var ParamError = zerror.New("A0100", "request params error")

func BindForm(ctx *gin.Context, form interface{}, verbose bool) bool {
	err := ctx.ShouldBind(form)
	if err != nil {
		if _errorLogger != nil {
			_errorLogger.Warnf("Bind request params failed, path: %s, form: %v, error: %v.", ctx.FullPath(), form, err)
		}
		if verbose {
			zresponse.Error(ctx, ParamError.Wrap(err))
		} else {
			zresponse.Error(ctx, ParamError)
		}
		return false
	}
	return true
}
