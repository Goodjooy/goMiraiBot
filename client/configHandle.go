package client

import (
	"goMiraiQQBot/constdata"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server   Server   `yaml:"server"`
	Bot      Bot      `yaml:"bot"`
	Database Database `yaml:"database"`
}

func (cfg Config) GetQQId() uint64 {
	return uint64(cfg.Bot.QQ)
}

type Server struct {
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
}

type Bot struct {
	QQ      uint64 `yaml:"QQ"`
	AuthKey string `yaml:"authKey"`
}
type Database struct {
	Eable bool `yaml:"enable"`

	Mode   constdata.DatabaseMode `yaml:"mode"`
	DbType string                 `yaml:"db"`

	DbName     string `yaml:"dbName"`
	DbUser     string `yaml:"dbUser"`
	DbPassword string `yaml:"dbPassword"`

	DbHost string `yaml:"dbHost"`
	DbPort uint   `yaml:"dbPort"`
}

func (cfg Config) IsEnable() bool {
	return cfg.Database.Eable
}
func (cfg Config) DbType() string {
	return cfg.Database.DbType
}
func (cfg Config) DbMode() constdata.DatabaseMode {
	return cfg.Database.Mode
}

func (cfg Config) GetDBUserName() string {
	return cfg.Database.DbUser
}
func (cfg Config) GetDBPassword() string {
	return cfg.Database.DbPassword
}
func (cfg Config) GetDBName() string {
	return cfg.Database.DbName
}

func (cfg Config) GetDBHost() string {
	return cfg.Database.DbHost
}
func (cfg Config) GetDBPort() uint {
	return cfg.Database.DbPort
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
		log.Printf("Open Config File Falure: %v| (Using Default Config)", err)
		return c
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read Config File Falure: %v| (Using Default Config)", err)
		return c
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Printf("UnMarshal(json) Config File Falure: %v| (Using Default Config)", err)
		return c
	}

	return cfg
}
