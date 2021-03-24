package router

import (
	"IntelligentTransfer/pkg/logger"
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
		logger.Errorf("Register user failed userInfo{%+v} err{%+v}", json, err)
		context.JSON(417, gin.H{"msg": "注册失败", "userId": "", "error": err.Error()})
	}
	context.JSON(http.StatusOK, gin.H{"msg": "登录成功", "userId": UUid})
}

//
