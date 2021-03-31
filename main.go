package main

import (
	"IntelligentTransfer/middleware"
	"IntelligentTransfer/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.Use(middleware.GinLogger())

	r.POST("/register", router.Register)
	r.POST("/login", router.Login)

	v1 := r.Group("/user")
	v1.Use(middleware.Cookie())
	{
		v1.POST("/:id/upload", router.Upload)
	}

	_ = r.Run(":40000")

}
