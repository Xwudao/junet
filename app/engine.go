package app

import (
	"github.com/gin-gonic/gin"
)

type Engine struct {
	*gin.Engine
}

func NewEngine(mode string) *Engine {
	gin.SetMode(mode)
	e := &Engine{Engine: gin.Default()}
	return e
}
