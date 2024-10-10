package config

import (
	"log"
	"os"

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
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Password string `json:"password"`
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
	data, err := os.ReadFile(configFile)

	if err != nil {
		log.Println(err)
		configFile = "config.json"
		data, err = os.ReadFile("/config/" + configFile)
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
