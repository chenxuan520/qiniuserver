package utils

import (
	"io/ioutil"
	"testing"
	"github.com/chenxuan520/qiniuserver/config"
)

func TestUploadPath(t *testing.T) {
	Init(config.GlobalConfig.Qiniu.AccessKey,config.GlobalConfig.Qiniu.SecretKey,config.GlobalConfig.Qiniu.Bucket)
	err := UploadFromPath("test/a.png", "/file/xiaozhu/file/download/sound.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUploadByte(t *testing.T) {
	Init(config.GlobalConfig.Qiniu.AccessKey,config.GlobalConfig.Qiniu.SecretKey,config.GlobalConfig.Qiniu.Bucket)
	data, err := ioutil.ReadFile("/file/xiaozhu/file/download/sound.png")
	if err != nil {
		t.Fatal(err)
	}

	err = UploadFromByte("test/b.png", data)
	if err != nil {
		t.Fatal(err)
	}
}
