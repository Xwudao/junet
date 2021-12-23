package app

import (
	"github.com/gin-gonic/gin"
)

type Engine struct {
	*gin.Engine
}

func NewEngine() *Engine {
	e := &Engine{Engine: gin.New()}
	return e
}
