package apiserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ResponseHandler(c *gin.Context, code int, msg string, data interface{}) {
	resp := map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	}
	c.JSON(200, resp)
}

func ginLogFormatFunc() func(arg gin.LogFormatterParams) string {

	f := func(arg gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %s %s %d %s  %s\n",
			arg.TimeStamp.Format("2006-01-02 15:04:05"),
			arg.ClientIP,
			arg.Request.Proto,
			arg.Method,
			arg.Path,
			arg.StatusCode,
			arg.Latency,
			arg.ErrorMessage,
		)
	}
	return f

}
