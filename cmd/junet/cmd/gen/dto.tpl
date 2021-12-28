{{- /*gotype: github.com/Xwudao/junet/cmd/junet/cmd/gen.DtoData*/ -}}
type {{.Name}}Dto struct {
{{range .Fields}}
    {{- .Name}} {{.Type}} `json:"{{.JsonName}}"`
{{end}}}