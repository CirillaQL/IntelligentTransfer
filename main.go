package main

import (
	"IntelligentTransfer/service"
	"fmt"
)

func main() {
	//r := gin.Default()
	//r.Use(middleware.Cors())
	//r.POST("/login", func(context *gin.Context) {
	//	json := make(map[string]interface{})
	//	_ = context.BindJSON(&json)
	//	logger.ZapLogger.Info(json["password"].(string))
	//})
	//_ = r.Run(":8090")
	//
	json := make(map[string]interface{})
	json["user_name"] = "dsd"
	json["nick_name"] = "cscd"
	json["sex"] = "男"
	json["province"] = "辽宁省"
	json["city"] = "大连市"
	json["address"] = "哦你的承诺"
	json["company"] = "腾讯"
	json["phone_number"] = "15840613358"
	json["email"] = "1194946223@qq.com"
	json["password"] = "ql1194946223"
	json["id_card"] = "210204199906135355"
	//service.Register(json)
	fmt.Println(service.LoginWithPassword("1194946223@qq.com", "ql1194946223", 2))
}
