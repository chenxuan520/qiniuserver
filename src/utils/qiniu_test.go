package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/chenxuan520/qiniuserver/config"
)

var (
	GlobalConfig *config.Config
)

func TestUploadPath(t *testing.T) {
	loadConfig()
	err := UploadFromPath("test/a.png", "/file/xiaozhu/file/download/sound.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetDirList(t *testing.T) {
	loadConfig()
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "0",
			args: args{
				dir: "test",
			},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDirList(tt.args.dir)
			t.Log(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDirList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func loadConfig() {
	configFile := "../../config/config.json"
	data, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Println(err)
		configFile = "config.json"
		data, err = ioutil.ReadFile("/config/" + configFile)
		if err != nil {
			log.Println("Read config error!")
			log.Panic(err)
			return
		}
	}

	config := &config.Config{}

	err = json.Unmarshal(data, config)

	if err != nil {
		log.Println("Unmarshal config error!")
		log.Panic(err)
		return
	}

	GlobalConfig = config
	log.Println("Config " + configFile + " loaded.")

	err = Init(GlobalConfig.Qiniu.AccessKey, GlobalConfig.Qiniu.SecretKey, GlobalConfig.Qiniu.Bucket, GlobalConfig.Qiniu.Zone)
	if err != nil {
		log.Println(err)
		log.Panic(err)
	}
}
