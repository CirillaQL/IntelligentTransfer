package service

import (
	"IntelligentTransfer/module"
	"IntelligentTransfer/pkg/encrypt"
	"IntelligentTransfer/pkg/logger"
	sql "IntelligentTransfer/pkg/mysql"
	"fmt"
	"strconv"
	"sync"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// OpenExcel 处理Excel
func OpenExcel(fileName string) []error {
	var errorList []error
	excel, err := excelize.OpenFile("./storage/" + fileName)
	if err != nil {
		logger.ZapLogger.Sugar().Errorf("Open Excel File{%+v} Failed, err:{%+v}", fileName, err)
		errorList = append(errorList, err)
		return errorList
	} else {
		logger.ZapLogger.Sugar().Debug("Open Excel " + fileName)
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

// GetMeeting 中读取到的Excel信息保存到DB-meetings
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
	meetingList := getMeetingsInfo(rows, SheetName)
	fmt.Println(meetingList)
	db := sql.GetDB()
	if db.Migrator().HasTable("meetings") {
		db.Table("meetings").Create(&meetingList)
	} else {
		db.Migrator().CreateTable(&module.Meeting{})
		db.Table("meetings").Create(&meetingList)
	}

}

// getMeetingsInfo 从Excel表中获取所有的数据并保存
func getMeetingsInfo(rows [][]string, meetingName string) []module.Meeting {
	var result []module.Meeting
	uuid := generateUUID()
	for index, meetingRow := range rows {
		if index == 0 {
			continue
		} else {
			if len(meetingRow) != 18 {
				break
			} else {
				MeetingInfo := module.Meeting{}
				MeetingInfo.UUid = generateUUID()
				MeetingInfo.MeetingUUid = uuid
				MeetingInfo.MeetingName = meetingName
				MeetingInfo.Name = meetingRow[0]
				MeetingInfo.Level = getUserLevel(meetingRow[1])
				MeetingInfo.Company = meetingRow[2]
				MeetingInfo.Sex = meetingRow[3]
				MeetingInfo.IdCard = meetingRow[4]
				MeetingInfo.PhoneNumber = meetingRow[5]
				MeetingInfo.IfOrderHotel = getIfOrder(meetingRow[6])
				MeetingInfo.IfOrderPlane = getIfOrder(meetingRow[7])
				MeetingInfo.StartDate = meetingRow[8]
				MeetingInfo.StartTime = meetingRow[9]
				MeetingInfo.StartBeginAddress = meetingRow[10]
				MeetingInfo.StartEndAddress = meetingRow[11]
				MeetingInfo.StartShift = meetingRow[12]
				MeetingInfo.ReturnDate = meetingRow[13]
				MeetingInfo.ReturnTime = meetingRow[14]
				MeetingInfo.ReturnStartAddress = meetingRow[15]
				MeetingInfo.ReturnEndAddress = meetingRow[16]
				MeetingInfo.ReturnShift = meetingRow[17]
				MeetingInfo.IfSolve = 0
				result = append(result, MeetingInfo)
			}
		}
	}
	return result
}

// GetMeetingExcel 从DB-meetings中读取数据保存到本地
func GetMeetingExcel(tableName, meetingUUid string) (string, error) {
	logger.ZapLogger.Sugar().Infof("GetMeetingInfo tableName: [%+v]  meetingUUid:[%+v]", tableName, meetingUUid)
	//根据表名获取对应的信息
	db := sql.GetDB()
	//首先获取接站的排序后信息
	var pickInfo []module.SmartMeeting
	db.Table(tableName).Order("pick_time").Where("meeting_u_uid = ? AND pick_or_sent = ?", meetingUUid, 1).Find(&pickInfo)
	logger.ZapLogger.Sugar().Info(pickInfo)
	//获取送站的相关信息
	var sentInfo []module.SmartMeeting
	db.Table(tableName).Order("sent_time").Where("meeting_u_uid = ? AND pick_or_sent = ?", meetingUUid, 0).Find(&sentInfo)
	logger.ZapLogger.Sugar().Info(sentInfo)
	//首先打开文件，设置sheet名称
	file := excelize.NewFile()
	file.SetSheetName("Sheet1", "接站")
	file.NewSheet("送站")
	//设置表头
	setSheetTitle(file)
	savePickInfo(file, pickInfo)
	saveSentInfo(file, sentInfo)
	//保存接站信息
	fileName := "./storage/" + tableName + ".xlsx"
	if err := file.SaveAs(fileName); err != nil {
		logger.ZapLogger.Sugar().Errorf("Save File Failed. Err:%+v", err)
	}
	return "", nil
}

//将输出表的表头信息保存到对应的值中
func setSheetTitle(file *excelize.File) {
	//设置表头的对应值
	file.SetCellValue("接站", "A1", "用户名")
	file.SetCellValue("接站", "B1", "电话号码")
	file.SetCellValue("接站", "C1", "到达站")
	file.SetCellValue("接站", "D1", "航班号")
	file.SetCellValue("接站", "E1", "时间")
	file.SetCellValue("接站", "F1", "司机名称")
	//设置表头的对应值
	file.SetCellValue("送站", "A1", "用户名")
	file.SetCellValue("送站", "B1", "电话号码")
	file.SetCellValue("送站", "C1", "到达站")
	file.SetCellValue("送站", "D1", "航班号")
	file.SetCellValue("送站", "E1", "时间")
	file.SetCellValue("送站", "F1", "司机名称")
}

//保存接站信息
func savePickInfo(file *excelize.File, smartMeetings []module.SmartMeeting) {
	i := 2
	for _, value := range smartMeetings {
		file.SetCellValue("接站", assembleCellString("A", i), value.UserName)
		file.SetCellValue("接站", assembleCellString("B", i), value.UserPhoneNumber)
		file.SetCellValue("接站", assembleCellString("C", i), value.ToAddress)
		file.SetCellValue("接站", assembleCellString("D", i), value.Shift)
		file.SetCellValue("接站", assembleCellString("E", i), value.PickTime)
		file.SetCellValue("接站", assembleCellString("F", i), getDriverInfo(value.DriverUUid))
		i++
	}
	//合并单元格
	for j := 0; j < len(smartMeetings)-1; j++ {
		if smartMeetings[j].DriverUUid == smartMeetings[j+1].DriverUUid {
			file.MergeCell("接站", assembleCellString("F", j+2), assembleCellString("F", j+3))
		}
	}
}

//保存送站信息
func saveSentInfo(file *excelize.File, smartMeetings []module.SmartMeeting) {
	i := 2
	for _, value := range smartMeetings {
		file.SetCellValue("送站", assembleCellString("A", i), value.UserName)
		file.SetCellValue("送站", assembleCellString("B", i), value.UserPhoneNumber)
		file.SetCellValue("送站", assembleCellString("C", i), value.ToAddress)
		file.SetCellValue("送站", assembleCellString("D", i), value.Shift)
		file.SetCellValue("送站", assembleCellString("E", i), value.SentTime)
		file.SetCellValue("送站", assembleCellString("F", i), getDriverInfo(value.DriverUUid))
		i++
	}
	//合并单元格
	for j := 0; j < len(smartMeetings)-1; j++ {
		if smartMeetings[j].DriverUUid == smartMeetings[j+1].DriverUUid {
			file.MergeCell("送站", assembleCellString("F", j+2), assembleCellString("F", j+3))
		}
	}
}

// 拼接字符和数字，生成单元格对应的字符串
func assembleCellString(alpha string, number int) string {
	return alpha + strconv.Itoa(number)
}

// 根据司机uuid获取司机信息
func getDriverInfo(driverUUid string) string {
	db := sql.GetDB()
	var driverName string
	db.Raw("select users.user_name from users inner join drivers on users.uuid ="+
		" drivers.user_u_uid where drivers.u_uid = ?", driverUUid).Scan(&driverName)
	var driverPhone string
	db.Raw("select users.phone_number from users inner join drivers on users.uuid ="+
		" drivers.user_u_uid where drivers.u_uid = ?", driverUUid).Scan(&driverPhone)
	if len(driverPhone) == 0 {
		return ""
	}
	realPhone, _ := encrypt.AesDecrypt(driverPhone)
	return driverName + "-" + realPhone
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

//判断是否订了飞机和酒店
func getIfOrder(input string) uint32 {
	if input == "0" {
		return 0
	} else if input == "1" {
		return 1
	} else {
		return 2
	}
}
