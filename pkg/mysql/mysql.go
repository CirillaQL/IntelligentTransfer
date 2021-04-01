package mysql

import (
	"IntelligentTransfer/config"
	log "IntelligentTransfer/pkg/logger"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

var db *gorm.DB
var once sync.Once

//初始化DB连接
func initDB() {
	once.Do(func() {
		var err error
		log.Info("begin init mysql db")
		mysqlConfig := config.GetConfig().Sub("mysql")
		user := mysqlConfig.GetString("user")
		password := mysqlConfig.GetString("password")
		url := mysqlConfig.GetString("host")
		dbName := mysqlConfig.GetString("db")
		connectSource := user + ":" + password + "@tcp" + "(" + url + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(connectSource), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Silent),
			SkipDefaultTransaction: true,
		})
		if err != nil {
			log.Errorf("init mysql db failed: %+v", err)
			panic(err)
		} else {
			log.ZapLogger.Sugar().Info("init mysql success")
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
