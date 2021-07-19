package http

import (
	"gift-code/internal/router"
	"github.com/gin-gonic/gin"
)

func InitRun()  {

	r:=gin.Default()
	//开启服务，监听请求
	router.Router(r)
	r.Run(":8080")

}
