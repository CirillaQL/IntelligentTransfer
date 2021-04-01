package middleware

import (
	"IntelligentTransfer/pkg/logger"
	"IntelligentTransfer/pkg/redis"
	"IntelligentTransfer/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Cors 中间件，处理前端axios发送的跨域请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "600")
		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

// GinLogger 中间件，将gin在控制台显示的log保存在本地的log文件上
func GinLogger() gin.HandlerFunc {
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
		logger.ZapLogger.Sugar().Infof("| %3d | %8v | %10s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

// Cookie 校验请求的Cookie
func Cookie() gin.HandlerFunc {
	return func(context *gin.Context) {
		var conn = redis.GetRedisConn()
		defer conn.Close()
		//获取Cookie
		getToken, err := context.Cookie("token")
		if err != nil {
			logger.ZapLogger.Sugar().Errorf("middleware get token failed Err:%+v", err)
			context.Abort()
		}
		value, err := redis.GetToken(getToken)
		if err != nil {
			logger.ZapLogger.Sugar().Errorf("middleware get token from redis failed Err:%+v", err)
			context.Abort()
		}
		tokenClaim, err := token.ParseToken(getToken)
		if err != nil {
			logger.ZapLogger.Sugar().Errorf("middleware parse token failed Err:%+v", err)
			context.Abort()
		}
		if tokenClaim.UUid == value {
			context.Next()
		} else {
			context.Abort()
		}
	}
}
