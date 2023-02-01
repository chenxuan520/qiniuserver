package utils

import (
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

var (
	upload  *storage.FormUploader  = nil
	manager *storage.BucketManager = nil
	cfg     *storage.Config        = nil
	zone    *storage.Zone          = nil
	mac     *auth.Credentials      = nil
)

func Init(access, secret, buc, area string) error {
	if len(access) == 0 || len(secret) == 0 || len(buc) == 0 {
		return fmt.Errorf("empty config")
	}

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

	cfg = &storage.Config{
		Zone:          zone,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	mac = auth.New(accessKey, secretKey)
	manager = storage.NewBucketManager(mac, cfg)
	upload = storage.NewFormUploader(cfg)

	return nil
}

func UploadFromPath(key, filepath string) error {
	if len(accessKey) == 0 || len(secretKey) == 0 || len(bucket) == 0 {
		return fmt.Errorf("empty config")
	}

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}

	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	upToken := putPolicy.UploadToken(mac)
	err := upload.PutFile(context.Background(), &ret, upToken, key, filepath, &putExtra)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(ret.Key, ret.Hash)
	return nil
}

func GetDirList(dir string) ([]string, error) {
	limit := 1000
	prefix := dir + "/"
	//初始列举marker为空
	marker := ""
	result := []string{}
	for {
		entries, _, nextMarker, hasNext, err := manager.ListFiles(bucket, prefix, "", marker, limit)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if len(entry.Key) <= len(prefix) {
				continue
			}
			key := entry.Key[len(prefix):]
			result = append(result, key)
		}
		if hasNext {
			marker = nextMarker
		} else {
			break
		}
	}
	return result, nil
}

func DeleteFile(path string) error {
	return manager.Delete(bucket, path)
}
