{{- /*gotype: github.com/Xwudao/junet/cmd/junet/cmd/gen.Route*/ -}}
package services

import (
)

type {{.RouteName}}Service struct {
}

func New{{.RouteName}}Service() *{{.RouteName}}Service {
	return &{{.RouteName}}Service{}
}

func (s *{{.RouteName}}Service) Index() {
}
