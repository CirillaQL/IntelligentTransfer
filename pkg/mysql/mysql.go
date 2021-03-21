package mysql

import (
	"IntelligentTransfer/config"
	"IntelligentTransfer/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"sync"
)

var db *gorm.DB
var once sync.Once

//初始化DB连接
func initDB() {
	once.Do(func() {
		var err error
		logger.Info("begin init mysql db")
		mysqlConfig := config.GetConfig().Sub("mysql")
		user := mysqlConfig.GetString("user")
		password := mysqlConfig.GetString("password")
		url := mysqlConfig.GetString("host")
		dbName := mysqlConfig.GetString("db")
		connectSource := user + ":" + password + "@" + "(" + url + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open("mysql", connectSource)
		if err != nil {
			logger.Errorf("init mysql db failed: %+v", err)
			panic(err)
		}
	})
}

//获取db实例
func GetDB() *gorm.DB {
	if db == nil {
		initDB()
	}
	return db
}
