package utils

import (
	"context"
	"fmt"
	"sort"

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

func GetDirList(dir string) ([]storage.ListItem, error) {
	limit := 1000
	prefix := dir
	if len(dir) > 0 {
		prefix = dir + "/"
	}
	//初始列举marker为空
	marker := ""
	var result []storage.ListItem
	result = make([]storage.ListItem, 0)
	for {
		entries, _, nextMarker, hasNext, err := manager.ListFiles(bucket, prefix, "", marker, limit)
		if err != nil {
			fmt.Printf("Error from manager.ListFiles for prefix '%s': %v\n", prefix, err)
			return nil, err
		}

		// sort by update time
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].PutTime > entries[j].PutTime
		})

		for i := range entries {
			if len(entries[i].Key) <= len(prefix) {
				continue
			}
			// Remove the prefix from the key, but keep the rest of the struct intact
			entries[i].Key = entries[i].Key[len(prefix):]
			result = append(result, entries[i])
		}

		if hasNext {
			marker = nextMarker
		} else {
			// no more items
			break
		}
	}
	return result, nil
}

func DeleteFile(path string) error {
	return manager.Delete(bucket, path)
}
