package router

import (
	"IntelligentTransfer/module"
	"IntelligentTransfer/pkg/encrypt"
	"IntelligentTransfer/pkg/logger"
	sql "IntelligentTransfer/pkg/mysql"
	"IntelligentTransfer/service"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"sync"
)

func Upload(c *gin.Context) {
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

func Download(c *gin.Context) {
	userId := c.Param("id")
	fileName := c.Param("name")
	//调用生成Excel文件的方法
	wg := sync.WaitGroup{}
	wg.Add(1)
	go createExcel(userId, fileName, &wg)
	wg.Wait()
	file := "./storage/" + fileName + ".xlsx"
	xlsx, _ := excelize.OpenFile(file)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName+".xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	_ = xlsx.Write(c.Writer)
}

//封装生成会议文件的方法，通过协程与WaitGroup同步
func createExcel(userId, fileName string, wg *sync.WaitGroup) {
	defer wg.Done()
	db := sql.GetDB()
	var user module.User
	db.Table("users").Where("uuid = ?", userId).Find(&user)
	user.PhoneNumber, _ = encrypt.AesDecrypt(user.PhoneNumber)
	logger.ZapLogger.Sugar().Infof("user %+v getMeetingExcel", user.UserName)
	var smartMeeting module.SmartMeeting
	db.Table(fileName).Where("user_phone_number = ?", user.PhoneNumber).Find(&smartMeeting)
	_, err := service.GetMeetingExcel(fileName, smartMeeting.MeetingUUid)
	if err != nil {
		return
	}
}
