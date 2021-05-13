package database

import "gorm.io/gorm"

var database *gorm.DB
var isAble bool = false

var dbModels []interface{}

func IsAbleUsingDataBase() bool {
	return isAble
}

func GetDB() *gorm.DB {
	return database
}
func GetDebugDB() *gorm.DB {
	return database.Debug()
}

func AsignDBModel(model interface{}) {
	dbModels = append(dbModels, model)
}


//TODO:数据库增删改查封装
