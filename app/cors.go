package app

import (
	"time"

	"github.com/gin-contrib/cors"
)

type CorsOpt func(*cors.Config)

func SetOriginFun(f func(origin string) bool) CorsOpt {
	return func(config *cors.Config) {
		config.AllowOriginFunc = f
	}
}
func SetMaxAge(age time.Duration) CorsOpt {
	return func(config *cors.Config) {
		config.MaxAge = age
	}
}
func SetCredentials(b bool) CorsOpt {
	return func(config *cors.Config) {
		config.AllowCredentials = b
	}
}
func SetExposeHeaders(s []string) CorsOpt {
	return func(config *cors.Config) {
		config.ExposeHeaders = s
	}
}
func SetHeaders(s []string) CorsOpt {
	return func(config *cors.Config) {
		config.AllowHeaders = s
	}
}
func SetMethods(s []string) CorsOpt {
	return func(config *cors.Config) {
		config.AllowMethods = s
	}
}
func SetOrigin(s []string) CorsOpt {
	return func(config *cors.Config) {
		config.AllowOrigins = s
	}
}
func (a *App) Cors(opts ...CorsOpt) {
	config := cors.DefaultConfig()
	for _, opt := range opts {
		opt(&config)
	}

	a.Use(cors.New(config))
}
