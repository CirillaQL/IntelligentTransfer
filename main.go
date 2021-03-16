package main

import (
	"IntelligentTransfer/pkg/log"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(log.GinLoggerToFile())

	router.Run()
}
