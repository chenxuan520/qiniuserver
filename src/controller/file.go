package controller

import (
	"io"
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

func UploadByUrl(g *gin.Context) {
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

	url := path.Path
	// 发起 GET 请求
	response, err := http.Get(url)
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}

	// 创建文件来保存下载内容
	fileName := utils.CreateFileName(config.GlobalConfig.Qiniu.FileName, utils.GenerateRandomString(6)+".url")
	file, err := os.Create(fileName)
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	// 将下载内容写入文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}
	defer response.Body.Close()

	// 保存到七牛云
	err = utils.UploadFromPath(config.GlobalConfig.Qiniu.UploadPath+"/"+fileName, fileName)
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}

	err = os.Remove(fileName)
	if err != nil {
		res.Error(g, http.StatusBadRequest, err.Error())
		return
	}

	res.Success(g, map[string]interface{}{
		"url": config.GlobalConfig.Qiniu.Domain + "/" + config.GlobalConfig.Qiniu.UploadPath + "/" + fileName,
	})
}

func UploadFile(g *gin.Context) {
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

	fileName := header.Filename
	if config.GlobalConfig.Qiniu.FileName != "" {
		fileName = utils.CreateFileName(config.GlobalConfig.Qiniu.FileName, fileName)
	}

	err = utils.UploadFromPath(config.GlobalConfig.Qiniu.UploadPath+"/"+fileName, header.Filename)
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
		"url": config.GlobalConfig.Qiniu.Domain + "/" + config.GlobalConfig.Qiniu.UploadPath + "/" + fileName,
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
