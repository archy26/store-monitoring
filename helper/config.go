package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var AppConfig *Config

type Config struct {
	Mysql       MysqlConfig `json:"mysql"`
	ServiceName string      `json:"serviceName"`
	Host        string      `json:"host"`
	Port        string      `json:"port"`
}

type MysqlConfig struct {
	Dsn string `json:"dsn"`
}

func Init() {
	f, err := os.Open("config/config.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	config := &Config{}
	err = json.NewDecoder(f).Decode(config)
	AppConfig = config
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	db, err := gorm.Open(mysql.Open(config.Mysql.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}
