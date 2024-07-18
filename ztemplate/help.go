package ztemplate

import (
	"html/template"
	"path/filepath"
)

var _must = template.Must

func base(filename string) string {
	return filepath.Base(filename)
}
