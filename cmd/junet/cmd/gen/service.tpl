{{- /*gotype: github.com/Xwudao/junet/cmd/junet/cmd/gen.Route*/ -}}
package services

import (
	"context"
	"gorm.io/gorm"
)

type {{.RouteName}}Service struct {
	db  *gorm.DB
	ctx context.Context
}

func New{{.RouteName}}Service(ctx context.Context) *{{.RouteName}}Service {
	return &{{.RouteName}}Service{
		//db:  db.GetDB(),
		ctx: ctx,
	}
}

func (s *{{.RouteName}}Service) Index() {
}
