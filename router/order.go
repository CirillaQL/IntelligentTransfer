package router

import (
	"IntelligentTransfer/module"
	"IntelligentTransfer/pkg/encrypt"
	"IntelligentTransfer/pkg/logger"
	sql "IntelligentTransfer/pkg/mysql"
	"IntelligentTransfer/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetUserOrders 查询订单路由
func GetUserOrders(content *gin.Context) {
	//获取用户id
	id := content.Param("id")
	//根据用户id查询用户信息
	db := sql.GetDB()
	var user module.User
	db.Where("uuid = ?", id).Find(&user)
	userPhone, err := encrypt.AesDecrypt(user.PhoneNumber)
	if err != nil {
		logger.ZapLogger.Sugar().Errorf("Decoding userPhone failed， err:%+v", err)
		content.JSON(http.StatusOK, gin.H{"msg": "失败"})
	}
	logger.ZapLogger.Sugar().Infof("user: %+v getOrders", user.UserName)
	orderList := service.GetOrders(userPhone)
	content.IndentedJSON(http.StatusOK, orderList)
}

// CancelOrder 删除订单路由
func CancelOrder(content *gin.Context) {
	//获取用户id
	id := content.Param("id")
	db := sql.GetDB()
	var user module.User
	db.Where("uuid = ?", id).Find(&user)
	//获取订单id
	uuid := content.Param("orderId")
	logger.ZapLogger.Sugar().Info(uuid)
	err := service.CancelUserOrder(uuid)
	if err != nil {
		logger.ZapLogger.Sugar().Errorf("User: %+v Cancel Order failed， err:%+v", user.UserName, err)
		content.JSON(http.StatusOK, gin.H{"msg": "1"})
	}
	content.JSON(http.StatusOK, gin.H{"msg": "0"})
}

// UpdateOrder 更新订单路由
func UpdateOrder(context *gin.Context) {
	json := make(map[string]interface{})
	err := context.BindJSON(&json)
	if err != nil {
		logger.ZapLogger.Sugar().Errorf("Json Bind Failed: %+v", err)
		context.JSON(http.StatusOK, gin.H{})
	}
	UUid := json["UUid"].(string)
	UserName := json["UserName"].(string)
	UserPhone := json["UserPhone"].(string)
	_ = service.UpdateOrder(UUid, UserName, UserPhone)
	context.JSON(http.StatusOK, gin.H{})
}

// GetDriverOrder 查找司机所有的订单
func GetDriverOrder(context *gin.Context) {
	id := context.Param("id")
	orderList := service.GetDriverOrder(id)
	context.IndentedJSON(http.StatusOK, orderList)
}
