package redis

import (
	"IntelligentTransfer/config"
	"IntelligentTransfer/pkg/logger"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

var pool *redis.Pool
var once sync.Once

// 初始化Redis池
func initRedis() {
	once.Do(func() {
		logger.ZapLogger.Sugar().Info("begin init redis")
		redisConfig := config.GetConfig().Sub("redis")
		pool = &redis.Pool{
			MaxIdle:     redisConfig.GetInt("maxIdle"),
			MaxActive:   redisConfig.GetInt("maxActive"),
			IdleTimeout: time.Duration(redisConfig.GetInt("idleTimeOut")) * time.Second,
			Wait:        false,
			Dial: func() (redis.Conn, error) {
				con, err := redis.Dial("tcp", redisConfig.GetString("address"),
					redis.DialConnectTimeout(time.Duration(redisConfig.GetInt64("conn_timeout"))*time.Second),
					redis.DialReadTimeout(time.Duration(redisConfig.GetInt64("read_timeout"))*time.Second),
					redis.DialWriteTimeout(time.Duration(redisConfig.GetInt64("write_timeout"))*time.Second),
				)
				if err != nil {
					logger.ZapLogger.Sugar().Errorf("redis init failed err: %+v", err)
					return nil, err
				}
				return con, nil
			},
		}
		logger.ZapLogger.Sugar().Info("init redis success")
	})
}

//获取redis连接
func GetRedisConn() redis.Conn {
	if pool == nil {
		initRedis()
	}
	return pool.Get()
}
