package database

import (
	"fmt"
	"goMiraiQQBot/lib/constdata"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type databaseConfig interface {
	IsEnable()bool
	DbType()string
	DbMode()constdata.DatabaseMode

	GetDBUserName() string
	GetDBPassword() string
	GetDBName() string

	GetDBHost() string
	GetDBPort() uint
}

type DBLinkerOpen func(dsn string) gorm.Dialector

func InitDatabaseConnect(dbOpen gorm.Dialector) *gorm.DB {

	
	db, err := gorm.Open(dbOpen, &gorm.Config{})
	if err != nil {
		log.Fatalf("Conntect To DataBase Fauilure : %v", err)
	}
	database = db
	isAble=true
	return db
}

func initMysqlConnect(dbCfg databaseConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%v:%v)/%s?charsetutf8&parseTime=True&loc=Local",
		dbCfg.GetDBUserName(),
		dbCfg.GetDBPassword(),
		dbCfg.GetDBHost(),
		dbCfg.GetDBPort(),
		dbCfg.GetDBName())
	return InitDatabaseConnect(mysql.Open(dsn))
}


func Init(dbCfg databaseConfig){
	if !dbCfg.IsEnable(){
		return
	}

	db:=initMysqlConnect(dbCfg)

	db.AutoMigrate(dbModels...)
}