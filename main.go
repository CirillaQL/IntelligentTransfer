package main

import (
	"IntelligentTransfer/middleware"
	"IntelligentTransfer/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

type Test struct {
	ID   int
	Name string
}

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
	r.POST("/login", router.Login)
	r.POST("/user/:id/upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "msg": fmt.Sprintf("error get form: %s", err.Error())})
			return
		}
		files := form.File["file"]
		for _, file := range files {
			basename := filepath.Base(file.Filename)
			filename := filepath.Join(".", basename)
			if err := c.SaveUploadedFile(file, filename); err != nil {
				c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "error": err.Error()})
				return
			}
		}

		var filenames []string
		for _, file := range files {
			filenames = append(filenames, file.Filename)
		}
		c.JSON(http.StatusOK, gin.H{"code": http.StatusAccepted, "msg": "upload ok!", "data": gin.H{"files": filenames}})
	})

	_ = r.Run(":40000")
	//测试gorm 自动建表
	//db := mysql.GetDB()
	//db.CreateTable(&module.Driver{})
	//db.DropTableIfExists(&Test{})
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
