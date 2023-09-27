package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var Gin *gin.Engine

type CustomController func(ctx *gin.Context) (Any, error)

type Any interface{}

type response struct {
	State   int	`json:"state"`
	Data    interface{}	`json:"data"`
	Msg     string 	`json:"msg"`
	TraceId string	`json:"traceId"`
}

// func init() {
// 	Gin = gin.New()
// 	Gin.Use(gin.Logger())
// 	Gin.Use(gin.Recovery())
// }

// 自定义 http handle
func ResponseConvert(c CustomController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := response{}
		_data, err := c(ctx)
		if err != nil {
			r.Msg = err.Error()
		}
		r.State = http.StatusOK
		r.Data = _data
		r.TraceId = ctx.MustGet("trace_id").(string)
		ctx.JSON(http.StatusOK, &r)
	}
}
