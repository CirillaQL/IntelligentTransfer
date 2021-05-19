package main

import (
	"IntelligentTransfer/middleware"
	"IntelligentTransfer/router"
	"IntelligentTransfer/service"
	_ "github.com/DeanThompson/ginpprof"
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
		v1.GET("/:id/getOrders", router.GetUserOrders)
		v1.GET("/:id/deleteOrder/:orderId", router.CancelOrder)
		v1.GET("/:id/getMeeting/:name", router.Download)
		v1.POST("/:id/checkMeeting", router.CheckUserMeetingInfo)
		v1.POST("/:id/updateMeeting", router.UpdateMeetingInfo)
	}
	_ = r.Run(":40000")

}
