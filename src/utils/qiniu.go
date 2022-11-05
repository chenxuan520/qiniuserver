package utils

import (
	// "bytes"
	"bytes"
	"context"
	"fmt"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

var (
	accessKey = ""
	secretKey = ""
	bucket    = ""
)

var zone *storage.Zone = nil

func Init(access, secret, buc, area string) error {
	accessKey = access
	secretKey = secret
	bucket = buc

	switch area {
	case "Huanan":
		zone = &storage.ZoneHuanan
	case "Huabei":
		zone = &storage.ZoneHuabei
	case "Huadong":
		zone = &storage.ZoneHuadong
	case "Xingjiapo":
		zone = &storage.ZoneXinjiapo
	default:
		return fmt.Errorf("wrong area config")
	}
	return nil
}

func UploadFromPath(key, filepath string) error {
	if len(accessKey) == 0 || len(secretKey) == 0 || len(bucket) == 0 {
		return fmt.Errorf("empty config")
	}

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := auth.New(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	cfg.Zone = zone
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	err := formUploader.PutFile(context.Background(), &ret, upToken, key, filepath, &putExtra)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(ret.Key, ret.Hash)
	return nil
}

func UploadFromByte(key string, data []byte, size int64) error {
	if len(accessKey) == 0 || len(secretKey) == 0 || len(bucket) == 0 {
		return fmt.Errorf("empty config")
	}

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := auth.New(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	//TODO change it to config
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	err := formUploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(data), size, &putExtra)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(ret.Key, ret.Hash)
	return nil
}
