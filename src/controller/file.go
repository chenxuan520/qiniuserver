package controller

import (
	"net/http"
	"os"

	"github.com/chenxuan520/qiniuserver/config"
	res "github.com/chenxuan520/qiniuserver/controller/response"
	"github.com/chenxuan520/qiniuserver/utils"
	"github.com/gin-gonic/gin"
)

func GetUploadFile(g *gin.Context) {
	//TODO add it to config
	header, err := g.FormFile("file")
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}

	err=g.SaveUploadedFile(header,header.Filename)
	if  err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}

	err=utils.Init(config.GlobalConfig.Qiniu.AccessKey, config.GlobalConfig.Qiniu.SecretKey, config.GlobalConfig.Qiniu.Bucket,config.GlobalConfig.Qiniu.Zone)
	if err != nil {
		res.Error(g,http.StatusBadRequest,err.Error())
		return
	}
	err = utils.UploadFromPath(config.GlobalConfig.Qiniu.UploadPath+"/"+header.Filename,header.Filename)
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}

	err=os.Remove(header.Filename);
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}

	res.Success(g, map[string]interface{}{
		"url": config.GlobalConfig.Qiniu.Domain + "/" + config.GlobalConfig.Qiniu.UploadPath + "/" + header.Filename,
	})
}
