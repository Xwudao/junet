package generator

import (
	_ "embed"
	"text/template"
)

//go:embed template.tpl
var qsCode string
var qsTmpl = template.Must(
	template.New("generator").
		Parse(qsCode),
)
