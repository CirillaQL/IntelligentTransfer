package service

import (
	"IntelligentTransfer/module"
	"IntelligentTransfer/pkg/logger"
	"IntelligentTransfer/pkg/mysql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"sync"
	"time"
)

func OpenExcel(fileName string) []error {
	var errorList []error
	excel, err := excelize.OpenFile("./storage/" + fileName)
	if err != nil {
		logger.ZapLogger.Sugar().Errorf("Open Excel File{%+v} Failed, err:{%+v}", fileName, err)
		errorList = append(errorList, err)
		return errorList
	}
	//获取会议名称：名称为Sheet表名
	MeetingNames := excel.GetSheetList()
	//并发从表中读取数据,并保存到DB
	wg := sync.WaitGroup{}
	errorChannel := make(chan error)
	for _, MeetingName := range MeetingNames {
		wg.Add(1)
		go GetMeeting(MeetingName, fileName, errorChannel, &wg)
	}
	go func() {
		wg.Wait()
		close(errorChannel)
	}()
	for val := range errorChannel {
		fmt.Println(val)
	}
	return nil
}

func GetMeeting(SheetName, fileName string, errorChannel chan error, wg *sync.WaitGroup) {
	//从每个Sheet表中读取对应的会议信息
	defer wg.Done()
	excel, err := excelize.OpenFile("./storage/" + fileName)
	if err != nil {
		logger.ZapLogger.Sugar().Errorf("goroutine can't open excel{%+v}. Error:{%+v}", fileName, err)
		errorChannel <- err
	}
	rows, err := excel.GetRows(SheetName)
	if err != nil {
		logger.ZapLogger.Sugar().Errorf("Open ExcelFile{%+v} ReadSheet{%+v} failed. Error: %+v",
			fileName, SheetName, err)
		errorChannel <- err
	}
	for index, row := range rows {
		if index == 0 {
			continue
		}
		err := SaveMeetingInfo(row, SheetName)
		if err != nil {
			errorChannel <- err
		}
	}
}

//从channel读取的值进行保存
func SaveMeetingInfo(meetingRow []string, meetingName string) error {
	//基础判断  模板表中有16列
	if len(meetingRow) != 16 {
		logger.ZapLogger.Sugar().Errorf("Excel SheetName:{%+v} has wrong column", meetingName)
		return fmt.Errorf("excel SheetName:{%+v} is not template", meetingName)
	}
	MeetingInfo := module.Meeting{}
	MeetingInfo.MeetingName = meetingName
	MeetingInfo.Name = meetingRow[0]
	MeetingInfo.Level = getUserLevel(meetingRow[1])
	MeetingInfo.Company = meetingRow[2]
	MeetingInfo.Sex = meetingRow[3]
	MeetingInfo.IdCard = meetingRow[4]
	MeetingInfo.PhoneNumber = meetingRow[5]
	MeetingInfo.IfOrderHotel = getIfOrder(meetingRow[6])
	MeetingInfo.IfOrderPlane = getIfOrder(meetingRow[7])
	ToTime, _ := time.ParseInLocation("2006-01-02 15:04", meetingRow[8], time.Local)
	MeetingInfo.ToTime = ToTime
	MeetingInfo.ToBeginAddress = meetingRow[9]
	MeetingInfo.ToEndAddress = meetingRow[10]
	MeetingInfo.ToShift = meetingRow[11]
	FromTime, _ := time.ParseInLocation("2006-01-02 15:04", meetingRow[12], time.Local)
	MeetingInfo.FromTime = FromTime
	MeetingInfo.FromBeginAddress = meetingRow[13]
	MeetingInfo.FromEndAddress = meetingRow[14]
	MeetingInfo.FromShift = meetingRow[15]
	db := mysql.GetDB()
	db.Create(&MeetingInfo)
	return nil
}

//获取用户对应的等级
func getUserLevel(input string) uint32 {
	if input == "组织者" {
		return 1
	} else if input == "讲师" {
		return 2
	} else if input == "参与人" {
		return 3
	} else {
		//身份未定
		return 4
	}
}

func getIfOrder(input string) uint32 {
	if input == "0" {
		return 0
	} else if input == "1" {
		return 1
	} else {
		return 2
	}
}
