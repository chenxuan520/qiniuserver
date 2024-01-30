package config

import (
	"io/ioutil"
	"log"

	"github.com/json-iterator/go"
)

type Qiniu struct {
	Domain     string `json:"domain"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Bucket     string `json:"bucket"`
	UploadPath string `json:"upload_path"`
	FileName   string `json:"file_name"`
	Zone       string `json:"zone"`
}

type Host struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type Config struct {
	Qiniu Qiniu `json:"qiniu"`
	Host  Host  `json:"host"`
}

var (
	GlobalConfig *Config
)

func Init() {
	configFile := "./config/config.json"
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

	config := &Config{}

	err = jsoniter.Unmarshal(data, config)

	if err != nil {
		log.Println("Unmarshal config error!")
		log.Panic(err)
		return
	}

	GlobalConfig = config
	log.Println("Config " + configFile + " loaded.")

}
