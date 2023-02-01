package main

import (
	"github.com/gin-gonic/gin"

	"github.com/chenxuan520/qiniuserver/config"
	"github.com/chenxuan520/qiniuserver/controller"
	"github.com/chenxuan520/qiniuserver/utils"
)

func main() {
	config.Init()
	err := utils.Init(config.GlobalConfig.Qiniu.AccessKey, config.GlobalConfig.Qiniu.SecretKey, config.GlobalConfig.Qiniu.Bucket, config.GlobalConfig.Qiniu.Zone)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(gin.Recovery())

	r.Static("/static", "./assert/")
	r.StaticFile("/", "./assert/index.html")

	r.GET("/info", controller.GetInfo)
	r.POST("/upload", controller.GetUploadFile)
	r.POST("/list", controller.GetFileList)
	r.POST("/delete", controller.DeleteFile)

	r.Run(config.GlobalConfig.Host.Ip + ":" + config.GlobalConfig.Host.Port)
}
