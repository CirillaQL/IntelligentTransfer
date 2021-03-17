package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

// 配置文件
var cfg *viper.Viper
var once sync.Once

// 初始化配置文件：本地读取config.yaml
func initCfg() {
	once.Do(func() {
		cfg = viper.New()
		cfg.SetConfigName("config")
		cfg.AddConfigPath("config/")
		cfg.SetConfigType("yaml")
		err := cfg.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("config getConfig failed"))
		}
	})
}

// 从文件读取配置
func GetConfig() *viper.Viper {
	if cfg == nil {
		initCfg()
	}
	return cfg
}
