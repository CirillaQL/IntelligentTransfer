package main

import (
	"IntelligentTransfer/middleware"
	"IntelligentTransfer/router"
	"github.com/gin-gonic/gin"
)

type Test struct {
	ID   int
	Name string
}

func main() {
	r := gin.Default()
	r.Use(middleware.Cors())

	r.POST("/register", router.Register)
	r.POST("/login", router.Login)
	r.POST("/user/:id/upload", router.Upload)

	_ = r.Run(":40000")
	//service.OpenExcel("测试会议.xlsx")

}
