package main

import (
	"IntelligentTransfer/middleware"
	"IntelligentTransfer/router"
	"IntelligentTransfer/service"

	"github.com/gin-gonic/gin"
)

func main() {
	go service.StartCron()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.GinLogger())

	r.POST("/register", router.Register)
	r.POST("/login", router.Login)

	v1 := r.Group("/user")
	v1.Use(middleware.Cookie())
	{
		v1.POST("/:id/upload", router.Upload)
		v1.POST("/:id/registerDriver", router.RegisterDriver)
	}

	_ = r.Run(":40000")
	//service.GetShiftInfo()
}
