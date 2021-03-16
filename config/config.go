package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

// 配置文件
var file *viper.Viper

var once sync.Once

// 初始化配置文件：本地读取config.yaml
func initCfg() {
	once.Do(func() {
		logger := logrus.WithFields(map[string]interface{}{
			"file": "config",
		})
		file = viper.New()
		file.SetConfigName("config") // name of config file (without extension)
		file.AddConfigPath("config/")
		file.SetConfigType("yaml")
		err := file.ReadInConfig()
		if err != nil {
			logger.Error(err)
		}
	})
}

// 从文件读取配置
func GetFile() *viper.Viper {
	if file == nil {
		initCfg()
	}
	return file
}

// 根据key获取根节点配置信息
func GetConfig(key string) interface{} {
	if file == nil {
		initCfg()
	}
	return file.Get(key)
}
