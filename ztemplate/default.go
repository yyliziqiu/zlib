package ztemplate

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

var GManager *Manager

func InitDefault(bases []string, htmls []string, funcs template.FuncMap) {
	GManager = New(bases, htmls, funcs)
}

func InitDefaultGlob(basePattern string, htmlPattern string, funcs template.FuncMap) {
	GManager = NewGlob(basePattern, htmlPattern, funcs)
}

func SetDebug(debug bool) *Manager {
	return GManager.SetDebug(debug)
}

func SetErrorTemplateName(name string) *Manager {
	return GManager.SetErrorTemplateName(name)
}

func Reload() *Manager {
	return GManager.Reload()
}

func Html(wr http.ResponseWriter, name string, data any) error {
	return GManager.Html(wr, name, data)
}

func HtmlGin(ctx *gin.Context, code int, name string, data any) {
	GManager.HtmlGin(ctx, code, name, data)
}

func PrintDefinedTemplates() {
	GManager.PrintDefinedTemplates()
}

func DefinedTemplates() []string {
	return GManager.DefinedTemplates()
}
