{{- /*gotype: github.com/Xwudao/junet/cmd/junet/cmd/gen.Route*/ -}}
package {{.PackageName}}

import (
	"github.com/Xwudao/junet/app"
	"github.com/gin-gonic/gin"
)

type {{.RouteName}}{{.RouteNameSuffix}} struct {
	Engine *app.App
}

func (r *{{.RouteName}}{{.RouteNameSuffix}}) SetUpRoutes() {
	{{.PackageName}} := r.Engine.Group("/{{.PackageName}}/{{.ToSnake .RouteName}}")
	{
		{{.PackageName}}.GET("", r.index())
	}
}

func (r *{{.RouteName}}{{.RouteNameSuffix}}) index() gin.HandlerFunc {
	return func(c *gin.Context) {
		app.SendJson(c, app.OkRtn("ok"))
	}
}
