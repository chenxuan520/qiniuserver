package main

import (
	"github.com/gin-gonic/gin"

	"github.com/chenxuan520/qiniuserver/config"
	"github.com/chenxuan520/qiniuserver/controller"
)

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())

	r.Static("/static", "./assert/")
	r.StaticFile("/", "./assert/index.html")

	r.POST("/upload", controller.GetUploadFile)

	r.Run(config.GlobalConfig.Host.Ip + ":" + config.GlobalConfig.Host.Port)
}
