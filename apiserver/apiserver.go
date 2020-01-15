package apiserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xloger/g"
)

func Start() {
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(ginLogFormatFunc()))

	v1 := r.Group("/agent/api")
	{
		v1.POST("/job/create", CreateCrontabAPI)
		v1.GET("/job/update", UpdateCrontabAPI)
		v1.GET("/job/delete", DeleteCrontabAPI)
	}

	sock := fmt.Sprintf("0.0.0.0:%d", g.Config().XlogNodeConfig.Port)
	r.Run(sock)

}
