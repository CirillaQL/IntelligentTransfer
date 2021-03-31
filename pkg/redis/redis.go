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

// StoreToken 保存Token进入redis, token值为登录时生成的Token值，过期时间默认为Token的过期时间2个小时
func StoreToken(token, userId string) error {
	conn := GetRedisConn()
	_, err := conn.Do("SET", token, userId)
	if err != nil {
		logger.ZapLogger.Sugar().Errorf("User: %+v store token into redis failed. err: %+v", userId, err)
		return err
	}
	_, err = conn.Do("EXPIRE", token, 2*3600)
	if err != nil {
		logger.ZapLogger.Sugar().Errorf("User: %+v store token into redis failed. err: %+v", userId, err)
		return err
	}
	return nil
}

// GetToken 从Redis中读取对应的信息
func GetToken(token string) (string, error) {
	conn := GetRedisConn()
	result, err := redis.String(conn.Do("Get", token))
	if err != nil {
		logger.ZapLogger.Sugar().Errorf("get Token from redis failed Err:%+v", err)
		return "", err
	}
	return result, nil
}
