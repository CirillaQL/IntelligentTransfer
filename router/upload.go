package router

import (
	"IntelligentTransfer/pkg/logger"
	"IntelligentTransfer/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func Upload(c *gin.Context) {
	//cookie, err := c.Cookie("token")
	//if err != nil {
	//	logger.ZapLogger.Sugar().Errorf("get User Cookie failed, err:%+v", err)
	//	c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "msg": fmt.Sprintf("error get form: %s",
	//		err.Error())})
	//	return
	//}
	//fmt.Println(cookie)
	form, err := c.MultipartForm()
	if err != nil {
		logger.ZapLogger.Sugar().Errorf("get upload file failed. err: {%+v}", err)
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "msg": fmt.Sprintf("error get form: %s",
			err.Error())})
		return
	}
	files := form.File["file"]
	for _, file := range files {
		basename := filepath.Base(file.Filename)
		filename := "./storage/" + basename
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "error": err.Error()})
			return
		}
		errList := service.OpenExcel(basename)
		if len(errList) == 0 {
			continue
		} else {
			fmt.Println(errList)
		}
	}
	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Filename)
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusAccepted, "msg": "upload ok!", "data": gin.H{"files": filenames}})
}