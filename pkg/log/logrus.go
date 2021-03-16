package log

import (
	"IntelligentTransfer/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"sync"
	"time"
)

var log *logrus.Logger
var once sync.Once

//初始化DB
func initLog() {
	once.Do(func() {
		log := logrus.New()
		//加载config.yaml
		logCfg := config.GetFile().Sub("log")
		logFilePath := fmt.Sprintf("config log: %s", logCfg.Get("logFilePath"))
		if dir, err := os.Getwd(); err == nil {
			logFilePath = dir + logFilePath
		}
		_ = os.MkdirAll(logFilePath, 0777)
		now := time.Now()
		logFileName := now.Format("2006-01-02") + ".log"
		//创建日志文件
		fileName := path.Join(logFilePath, logFileName)
		if _, err := os.Stat(fileName); err != nil {
			_, _ = os.Create(fileName)
		}
		//写入文件
		src, _ := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		//设置输出
		log.Out = src
		//设置日志级别
		log.SetLevel(logrus.DebugLevel)
		//设置日志格式
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
		log.Debugf("pkg log init success")
	})
}

//获取Log实例
func Logger() *logrus.Logger {
	if log == nil {
		initLog()
	}
	return log
}

//将日志文件保存到本地文件
func GinLoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		//日志格式
		log.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
