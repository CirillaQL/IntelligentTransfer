package router

import (
	"IntelligentTransfer/pkg/logger"
	"IntelligentTransfer/pkg/redis"
	"IntelligentTransfer/pkg/token"
	"IntelligentTransfer/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//注册路由
func Register(context *gin.Context) {
	json := make(map[string]interface{})
	_ = context.BindJSON(&json)
	UUid, err := service.Register(json)
	if err != nil {
		logger.ZapLogger.Sugar().Errorf("Register user failed userInfo{%+v} err{%+v}", json, err)
		context.JSON(http.StatusOK, gin.H{"msg": "注册失败", "userId": "", "error": err.Error()})
	} else {
		tokenString, err := token.GenToken(UUid, json["phoneNumber"].(string))
		if err != nil {
			logger.ZapLogger.Sugar().Errorf("Register user success but genToken failed. err:%+v", err)
			context.JSON(http.StatusOK, gin.H{"msg": "注册失败", "userId": UUid, "token": tokenString})
		}
		err = redis.StoreToken(tokenString, UUid)
		if err != nil {
			logger.ZapLogger.Sugar().Errorf("User: %+v storeToken Failed. Err: %+v", UUid, err)
			context.JSON(http.StatusOK, gin.H{"msg": "注册失败", "userId": UUid, "token": tokenString})
		}
		context.JSON(http.StatusOK, gin.H{"msg": "注册成功", "userId": UUid, "token": tokenString})
	}
}

//登录路由
func Login(context *gin.Context) {
	json := make(map[string]interface{})
	_ = context.Bind(&json)
	userInfo, password, inputType := loginJson(json)
	result, userId, phoneNumber, err := service.LoginWithPassword(userInfo, password, inputType)
	if err != nil || result == false {
		logger.ZapLogger.Sugar().Errorf("user login failed err:%+v, result:%+v", err, result)
		context.JSON(http.StatusOK, gin.H{"msg": "登录失败"})
	} else {
		logger.ZapLogger.Sugar().Info("user login success")
		tokenString, err := token.GenToken(userId, phoneNumber)
		if err != nil {
			logger.ZapLogger.Sugar().Errorf("User: %+v genToken Failed. Err: %+v", userId, err)
			context.JSON(http.StatusOK, gin.H{"msg": "登录失败"})
		}
		err = redis.StoreToken(tokenString, userId)
		if err != nil {
			logger.ZapLogger.Sugar().Errorf("User: %+v storeToken Failed. Err: %+v", userId, err)
			context.JSON(http.StatusOK, gin.H{"msg": "登录失败"})
		}
		context.JSON(http.StatusOK, gin.H{"msg": "登录成功", "token": tokenString, "userId": userId})
	}
}

//登录时校验json
func loginJson(json map[string]interface{}) (string, string, uint32) {
	if json["loginType"].(float64) == 1 {
		//输入为邮箱
		return json["userinfo"].(string), json["password"].(string), 1
	} else if json["loginType"].(float64) == 2 {
		//输入为邮箱
		return json["userinfo"].(string), json["password"].(string), 2
	} else {
		return "", "", 0
	}
}

func RegisterDriver(context *gin.Context) {
	json := make(map[string]interface{})
	_ = context.BindJSON(&json)
	s := service.DriverRegister(json)
	fmt.Println(s)
}
