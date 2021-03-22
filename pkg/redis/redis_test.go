package redis

import (
	"github.com/gomodule/redigo/redis"
	"testing"
	"time"
)

func TestGetRedisConn(t *testing.T) {
	pool = &redis.Pool{
		MaxIdle:     5,
		MaxActive:   0,
		IdleTimeout: time.Duration(10) * time.Second,
		Wait:        false,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				//logger.Errorf("redis init failed err: %+v", err)
				return nil, err
			}
			return con, nil
		},
	}
	conn := pool.Get()
	conn.Do("SET", "sdad", "dsd")
}
