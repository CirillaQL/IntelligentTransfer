package router

import (
	"IntelligentTransfer/pkg/logger"
	"IntelligentTransfer/pkg/token"
	"IntelligentTransfer/service"
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
		context.JSON(http.StatusOK, gin.H{"msg": "登录成功", "userId": UUid})
	}
}

//登录路由
func Login(context *gin.Context) {
	json := make(map[string]interface{})
	_ = context.Bind(&json)
	userInfo, password, inputType := loginJson(json)
	result, userId, err := service.LoginWithPassword(userInfo, password, inputType)
	if err != nil || result == false {
		logger.ZapLogger.Sugar().Errorf("user login failed err:%+v", err)
		context.JSON(http.StatusOK, gin.H{"msg": "登录失败"})
	} else {
		logger.ZapLogger.Sugar().Info("user login success")
		tokenString, _ := token.GenToken(userInfo, password)
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
