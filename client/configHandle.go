package client

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server Server `yaml:"server"`
	Bot    Bot    `yaml:"bot"`
}

func (cfg Config)GetQQId()uint64{
	return uint64(cfg.Bot.QQ)
}

type Server struct {
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
}

type Bot struct {
	QQ      uint64   `yaml:"QQ"`
	AuthKey string `yaml:"authKey"`
}

func LoadConfig(configPath string) Config {

	var path string
	if configPath == "" {
		path = "./config.yml"
	}
	path = configPath

	log.Printf("loading Config File[%v]", path)

	c := Config{
		Server: Server{
			Host: "127.0.0.1",
			Port: 8080,
		},
		Bot: Bot{
			QQ:      0,
			AuthKey: "",
		},
	}
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Open Config File Falure: %v| (Using Default Config)", err)
		return c
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Read Config File Falure: %v| (Using Default Config)", err)
		return c
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("UnMarshal(json) Config File Falure: %v| (Using Default Config)", err)
		return c
	}

	return cfg
}
