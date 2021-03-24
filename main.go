package main

import (
	"IntelligentTransfer/middleware"
	"IntelligentTransfer/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(middleware.Cors())
	//r.POST("/register", func(context *gin.Context) {
	//	json := make(map[string]string)
	//	_ = context.BindJSON(&json)
	//	fmt.Println(json["username"], json["password"], 2)
	//	//fmt.Println(json["username"].(string), json["password"].(string), 2)
	//	context.JSON(http.StatusOK, gin.H{"msg": "登录成功", "data": "OK"})
	//})
	r.POST("/register", router.Register)
	_ = r.Run(":40000")
	//
	//json := make(map[string]interface{})
	//json["user_name"] = "dsd"
	//json["nick_name"] = "cscd"
	//json["sex"] = "男"
	//json["province"] = "辽宁省"
	//json["city"] = "大连市"
	//json["address"] = "哦你的承诺"
	//json["company"] = "腾讯"
	//json["phone_number"] = "15840613358"
	//json["email"] = "1194946223@qq.com"
	//json["password"] = "ql1194946223"
	//json["id_card"] = "210204199906135355"
	////service.Register(json)
	//fmt.Println(service.LoginWithPassword("1194946223@qq.com", "ql1194946223", 2))
}
