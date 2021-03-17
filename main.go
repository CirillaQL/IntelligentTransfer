package main

import (
	"IntelligentTransfer/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(middleware.GinLogger())
	router.Run()
}
