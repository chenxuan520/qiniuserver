package controller

import (
	"net/http"
	"os"

	"github.com/chenxuan520/qiniuserver/config"
	res "github.com/chenxuan520/qiniuserver/controller/response"
	"github.com/chenxuan520/qiniuserver/utils"
	"github.com/gin-gonic/gin"
)

type ReqPath struct {
	Path string `json:"path"`
}

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

	err = g.SaveUploadedFile(header, header.Filename)
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}

	err = utils.UploadFromPath(config.GlobalConfig.Qiniu.UploadPath+"/"+header.Filename, header.Filename)
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}

	err = os.Remove(header.Filename)
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}

	res.Success(g, map[string]interface{}{
		"url": config.GlobalConfig.Qiniu.Domain + "/" + config.GlobalConfig.Qiniu.UploadPath + "/" + header.Filename,
	})
}

func GetFileList(g *gin.Context) {
	path := &ReqPath{}
	err := g.Bind(path)
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}
	if path.Path == "" {
		res.Error(g, http.StatusBadRequest, "cannot get dir")
		return
	}

	list, err := utils.GetDirList(path.Path)
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}

	res.Success(g, map[string]interface{}{
		"list": list,
	})
}

func DeleteFile(g *gin.Context) {
	path := &ReqPath{}
	err := g.Bind(path)
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}
	if path.Path == "" {
		res.Error(g, http.StatusBadRequest, "cannot get dir")
		return
	}

	err = utils.DeleteFile(path.Path)
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}

	res.Success(g, nil)
}

func GetInfo(g *gin.Context) {
	res.Success(g, map[string]interface{}{
		"path":   config.GlobalConfig.Qiniu.UploadPath,
		"domain": config.GlobalConfig.Qiniu.Domain,
	})
}
