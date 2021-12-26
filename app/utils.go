package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Xwudao/junet"
)

func CmnRtn(code int, msg string, v interface{}) gin.H {
	return gin.H{
		"code": code,
		"msg":  msg,
		"data": v,
	}
}
func OkRtn(v interface{}) gin.H {
	return CmnRtn(200, junet.SuccessRtn, v)
}
func ErrRtn(msg string) gin.H {
	return CmnRtn(0, msg, nil)
}
func SendJson(c *gin.Context, v interface{}) {
	c.JSON(http.StatusOK, v)
}
