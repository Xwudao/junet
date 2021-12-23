package app

import (
	"github.com/gin-gonic/gin"
)

type Context struct {
	context *gin.Context
}

func NewContext(context *gin.Context) *Context {
	return &Context{context: context}
}
