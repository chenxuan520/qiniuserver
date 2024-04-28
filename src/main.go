package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/chenxuan520/qiniuserver/config"
	"github.com/chenxuan520/qiniuserver/controller"
	"github.com/chenxuan520/qiniuserver/middlerware"
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

	api := r.Group("/api")
	api.Use(middlerware.PasswdAuth())
	{
		api.GET("/info", controller.GetInfo)
		api.POST("/upload", controller.UploadFile)
		api.POST("/upload_url", controller.UploadByUrl)
		api.POST("/list", controller.GetFileList)
		api.POST("/delete", controller.DeleteFile)
		// this api for test
		api.POST("/upload_test", controller.UploadTest)
	}

	log.Println("Runner in http://" + config.GlobalConfig.Host.Ip + ":" + config.GlobalConfig.Host.Port)
	r.Run(config.GlobalConfig.Host.Ip + ":" + config.GlobalConfig.Host.Port)
}
