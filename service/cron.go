package service

import (
	"IntelligentTransfer/constant"
	errorInfo "IntelligentTransfer/error"
	"IntelligentTransfer/module"
	"IntelligentTransfer/pkg/logger"
	sql "IntelligentTransfer/pkg/mysql"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

var c *cron.Cron

//定时任务初始化
func init() {
	c = cron.New()
	logger.ZapLogger.Sugar().Info("Cron init success")
}

// StartCron 开启定时任务
func StartCron() {
	logger.ZapLogger.Sugar().Info("Start Cron service success")
	c.AddFunc("@every 10s", CreateTable)
	c.AddFunc("@every 10s", GetTodayPickInfoByOrder)

	//c.AddFunc("@every 1m", testUpdate)
	c.Start()
	select {}
}

// CreateTable 生成以日期为TableName的数据库表
func CreateTable() {
	db := sql.GetDB()
	var result []module.Meeting
	db.Where("if_solve = ?", 0).Find(&result)
	ans := getDateMap(&result)
	if len(result) == 0 {
		return
	} else {
		logger.ZapLogger.Sugar().Info("Cron work begin to solve data")
		//根据生成的map结构保存数据到DB中去
		for date, meetings := range *ans {
			if db.Migrator().HasTable(date) == false {
				err := db.Table(date).Migrator().CreateTable(module.SmartMeeting{})
				if err != nil {
					logger.ZapLogger.Sugar().Errorf("Cron Create Table{%+v} failed. err: %+v", date, err)
					continue
				}
			}
			/*
				将key-date中的所有元素保存到DB中， 首先将Meeting信息拼装成SmartMeeting信息
			*/
			for _, meeting := range meetings {
				//开完会送站
				if meeting.IfPick == 0 {
					smartmeeting := module.SmartMeeting{}
					smartmeeting.UUid = generateUUID()
					smartmeeting.MeetingUUid = meeting.MeetingInfo.MeetingUUid
					smartmeeting.UserName = meeting.MeetingInfo.Name
					smartmeeting.UserPhoneNumber = meeting.MeetingInfo.PhoneNumber
					smartmeeting.LeveL = meeting.MeetingInfo.Level
					smartmeeting.FromAddress = meeting.MeetingInfo.ReturnStartAddress
					smartmeeting.ToAddress = meeting.MeetingInfo.ReturnEndAddress
					smartmeeting.SentTime = meeting.MeetingInfo.ReturnTime
					smartmeeting.Shift = meeting.MeetingInfo.ReturnShift
					smartmeeting.PickOrSent = 0
					smartmeeting.DriverUUid = ""
					db.Table(date).Create(&smartmeeting)
				}
				//准备开会，接站
				if meeting.IfPick == 1 {
					smartmeeting := module.SmartMeeting{}
					smartmeeting.UUid = generateUUID()
					smartmeeting.MeetingUUid = meeting.MeetingInfo.MeetingUUid
					smartmeeting.UserName = meeting.MeetingInfo.Name
					smartmeeting.UserPhoneNumber = meeting.MeetingInfo.PhoneNumber
					smartmeeting.LeveL = meeting.MeetingInfo.Level
					smartmeeting.FromAddress = meeting.MeetingInfo.StartBeginAddress
					smartmeeting.ToAddress = meeting.MeetingInfo.StartEndAddress
					smartmeeting.PickTime = meeting.MeetingInfo.StartTime
					smartmeeting.Shift = meeting.MeetingInfo.StartShift
					smartmeeting.PickOrSent = 1
					smartmeeting.DriverUUid = ""
					db.Table(date).Create(&smartmeeting)
				}
			}
		}
		//将if_solve的值更新为1
		db.Model(&module.Meeting{}).Where("if_solve = ?", 0).Update("if_solve", 1)
	}
}

//获取当前日期
func getToday() string {
	date := time.Now().Format("2006-01-02")
	return date
}

// GetTodayPickInfoByOrder 按照时间从表中获取数据，进行排序，目前获取数据仅仅为获取到当天的数据，其他之后时间段的暂时并不需要
func GetTodayPickInfoByOrder() {
	db := sql.GetDB()
	if db.Migrator().HasTable("2021-05-14") {
		go dueTodayPick("2021-05-14")
		go dueTodaySent("2021-05-14")
	}
	go GeneOrder("2021-05-14")
}

//封装获取到的接站信息
func dueTodayPick(tableName string) {
	smartMeetingPick := make([]module.SmartMeeting, 0)
	db := sql.GetDB()
	//此处获取的为接站
	db.Table(tableName).Order("pick_time").Where("pick_or_sent = ? AND driver_u_uid = ?", 1, "").Find(&smartMeetingPick)
	if len(smartMeetingPick) == 0 {
		return
	}
	//获取到按照时间进行了排序的接站信息表
	pickTimeMap := createPickTimeMap(smartMeetingPick)
	pickResult := assignmentDrivers(*pickTimeMap)
	for _, value := range *pickResult {
		for _, v := range value {
			if v.DriverUUid != "" {
				logger.ZapLogger.Sugar().Infof("user: %+v has Driver: %+v", v.UserName, v.DriverUUid)
				_ = updateDriverInfoToDB(tableName, v.UUid, v.DriverUUid)
			} else {
				//logger.ZapLogger.Sugar().Infof("user: %+v doesn't have driver! ", v.UserName)
				// TODO:后续会给会议组织者进行提醒
				continue
			}
		}
	}
}

//封装获取到的送站信息
func dueTodaySent(tableName string) {
	smartMeetingSent := make([]module.SmartMeeting, 0)
	db := sql.GetDB()
	//此时获取的为送站
	db.Table(tableName).Order("sent_time").Where("pick_or_sent = ? AND driver_u_uid = ?", 0, "").Find(&smartMeetingSent)
	if len(smartMeetingSent) == 0 {
		return
	}
	//获取到按照时间进行了排序的接站信息表
	sentTimeMap := createSentTimeMap(smartMeetingSent)
	sentResult := assignmentDrivers(*sentTimeMap)
	for _, value := range *sentResult {
		for _, v := range value {
			if v.DriverUUid != "" {
				logger.ZapLogger.Sugar().Infof("user: %+v has Driver: %+v", v.UserName, v.DriverUUid)
				_ = updateDriverInfoToDB(tableName, v.UUid, v.DriverUUid)
			} else {
				//logger.ZapLogger.Sugar().Infof("user: %+v doesn't have driver! ", v.UserName)
				// TODO:后续会给会议组织者进行提醒
				continue
			}
		}
	}
}

// 根据排序后获取的user表生成按照送站时间生成的map表
func createSentTimeMap(users []module.SmartMeeting) *map[string][]module.SmartMeeting {
	var result map[string][]module.SmartMeeting
	result = make(map[string][]module.SmartMeeting)
	for _, v := range users {
		if result[v.SentTime] == nil {
			var smartMeetings []module.SmartMeeting
			smartMeetings = append(smartMeetings, v)
			result[v.SentTime] = smartMeetings
		} else {
			result[v.SentTime] = append(result[v.SentTime], v)
		}
	}
	return &result
}

// 根据排序后获取的user表生成按照接站时间生成的map表
func createPickTimeMap(users []module.SmartMeeting) *map[string][]module.SmartMeeting {
	var result map[string][]module.SmartMeeting
	result = make(map[string][]module.SmartMeeting)
	for _, v := range users {
		if result[v.PickTime] == nil {
			var smartMeetings []module.SmartMeeting
			smartMeetings = append(smartMeetings, v)
			result[v.PickTime] = smartMeetings
		} else {
			result[v.PickTime] = append(result[v.PickTime], v)
		}
	}
	return &result
}

// assignmentDrivers 根据司机情况分配司机
func assignmentDrivers(users map[string][]module.SmartMeeting) *map[string][]module.SmartMeeting {
	for timeToGet, userList := range users {
		/*
			此处采用尽量少的车辆策略，根据某个时间段到达的人数不同，采用不同的分配策略
		*/
		numberOfPassenger := len(userList)
		//当该时间的用户数小于等于2人，则优先安排小轿车
		if numberOfPassenger <= 2 {
			smallCar := GetAllTypeOneDriver()
			afterAssignment := assignmentSmallCar(userList, smallCar)
			users[timeToGet] = afterAssignment
		} else if 3 <= numberOfPassenger && 5 >= numberOfPassenger {
			//该时间的用户数为3-5人，优先安排SUV
			suv := GetAllTypeTwoDriver()
			afterAssignment := assignmentSuv(userList, suv)
			users[timeToGet] = afterAssignment
		} else if 6 <= numberOfPassenger && 13 >= numberOfPassenger {
			//该时间的用户数为6-13人，优先安排考斯特
			coaster := GetAllTypeThreeDriver()
			afterAssignment := assignmentCoaster(userList, coaster)
			users[timeToGet] = afterAssignment
		} else if 14 <= numberOfPassenger {
			//该时间的用户数为14人，优先安排大巴车
			bus := GetAllTypeFourDriver()
			afterAssignment := assignmentBus(userList, bus)
			users[timeToGet] = afterAssignment
		}
	}
	return &users
}

// 分配司机后，将分配后的司机状态更新到DB中
func updateDriverInfoToDB(tableName, smartMeetingUUid, driverUUid string) error {
	db := sql.GetDB()
	if err := db.Table(tableName).Where("u_uid = ?", smartMeetingUUid).Update("driver_u_uid", driverUUid).Error; err != nil {
		logger.ZapLogger.Sugar().Errorf("SmartMeeting: %+v Update Driver Failed. Err: %+v", smartMeetingUUid, err)
		return err
	}
	return nil
}

/*
	此函数为分配小轿车的算法，其中小脚测的载客数为1-2人，思想为：
		1.如果用户数量少于一辆小轿车的最大载客量2人，那么直接分配一辆就可以
		2.如果用户数量多于一辆小轿车的最大载客量2人，那么求余数，如果余数为0，说明可以按照最大载客量全部坐满，
		  如果余数不为0，那么将最大载客量的整数倍装满后，剩下的余数个再去装
*/
func assignmentSmallCar(users []module.SmartMeeting, drivers []module.Driver) []module.SmartMeeting {
	if len(drivers) < 1 {
		return users
	}
	numOfUsers := len(users)
	if numOfUsers == 1 {
		users[0].DriverUUid = drivers[0].UUid
		UpdateDriverType(drivers[0].UUid, constant.DRIVER_WORKING)
	} else if numOfUsers == 2 {
		users[0].DriverUUid = drivers[0].UUid
		users[1].DriverUUid = drivers[0].UUid
		UpdateDriverType(drivers[0].UUid, constant.DRIVER_WORKING)
	} else {
		if numOfUsers%2 == 0 {
			j := 0
			for i := 0; i < numOfUsers-1; i = i + 2 {
				users[i].DriverUUid = drivers[j].UUid
				users[i+1].DriverUUid = drivers[j].UUid
				UpdateDriverType(drivers[j].UUid, constant.DRIVER_WORKING)
				j++
				if j == len(drivers) {
					break
				}
			}
		} else {
			j := 0
			for i := 0; i < numOfUsers-1; i = i + 2 {
				users[i].DriverUUid = drivers[j].UUid
				users[i+1].DriverUUid = drivers[j].UUid
				UpdateDriverType(drivers[j].UUid, constant.DRIVER_WORKING)
				j++
				if j == len(drivers) {
					break
				}
			}
			if j <= len(drivers) {
				users[numOfUsers-1].DriverUUid = drivers[j].UUid
				UpdateDriverType(drivers[j].UUid, constant.DRIVER_WORKING)
			}
		}
	}
	return users
}

/*
	此函数为分配别克SUV的算法，其中SUV的载客数为3-5人，思想为：
		1.如果用户数量少于一辆SUV的最大载客量5人，那么直接分配一辆就可以
		2.如果用户数量多于一辆SUV的最大载客量5人，那么求余数，如果余数为0，说明可以按照最大载客量全部坐满，
		  如果余数不为0，那么将最大载客量的整数倍装满后，剩下的余数个再去装
*/
func assignmentSuv(users []module.SmartMeeting, drivers []module.Driver) []module.SmartMeeting {
	if len(drivers) < 1 {
		return users
	}
	numOfUsers := len(users)
	//此时待安排的用户数少于5人，因此一辆SUV便可以装下所有人
	if numOfUsers <= 5 {
		for i := 0; i < numOfUsers; i++ {
			users[i].DriverUUid = drivers[0].UUid
			UpdateDriverType(drivers[0].UUid, constant.DRIVER_WORKING)
		}
		return users
	} else {
		//此时待安排的用户数多于5，需要多辆SUV
		//余数为0，正好全部坐满
		if numOfUsers%5 == 0 {
			j := 0
			for i := 0; i < numOfUsers-4; i = i + 5 {
				for k := i; k < 5; k++ {
					users[k].DriverUUid = drivers[j].UUid
					UpdateDriverType(drivers[j].UUid, constant.DRIVER_WORKING)
				}
				j++
				if j == len(drivers) {
					break
				}
			}
		} else {
			//余数为其他数
			numberOfLeftUser := numOfUsers % 5
			j := 0
			for i := 0; i <= numOfUsers-numberOfLeftUser; i = i + 5 {
				for k := i; k < i+5; k++ {
					users[k].DriverUUid = drivers[j].UUid
					UpdateDriverType(drivers[j].UUid, constant.DRIVER_WORKING)
				}
				j++
				if j == len(drivers) {
					break
				}
			}
			if len(drivers) > j {
				//此时还有剩余车辆，用来装剩余的人
				for i := numOfUsers - 1; i >= numOfUsers-numberOfLeftUser; i-- {
					users[i].DriverUUid = drivers[j+1].UUid
					UpdateDriverType(drivers[j+1].UUid, constant.DRIVER_WORKING)
				}
			}
		}
		return users
	}
}

/*
	此函数为分配考斯特Coaster的算法，其中考斯特的载客数为6-13人，思想为：
		1.如果用户数量少于一辆考斯特Coaster的最大载客量13人，那么直接分配一辆就可以
		2.如果用户数量多于一辆考斯特Coaster的最大载客量13人，那么求余数，如果余数为0，说明可以按照最大载客量全部坐满，
		  如果余数不为0，那么将最大载客量的整数倍装满后，剩下的余数个再去装
*/
func assignmentCoaster(users []module.SmartMeeting, drivers []module.Driver) []module.SmartMeeting {
	if len(drivers) < 1 {
		return users
	}
	numOfUsers := len(users)
	//此时待安排的用户数少于13人，因此一辆考斯特便可以装下所有人
	if numOfUsers <= 13 {
		for i := 0; i < numOfUsers; i++ {
			users[i].DriverUUid = drivers[0].UUid
			UpdateDriverType(drivers[0].UUid, constant.DRIVER_WORKING)
		}
		return users
	} else {
		//此时待安排的用户数多于13，需要多辆考斯特
		//余数为0，正好全部坐满
		if numOfUsers%13 == 0 {
			j := 0
			for i := 0; i < numOfUsers-12; i = i + 13 {
				for k := i; k < 13; k++ {
					users[k].DriverUUid = drivers[j].UUid
					UpdateDriverType(drivers[j].UUid, constant.DRIVER_WORKING)
				}
				j++
				if j == len(drivers) {
					break
				}
			}
		} else {
			//余数为其他数
			numberOfLeftUser := numOfUsers % 13
			j := 0
			for i := 0; i <= numOfUsers-numberOfLeftUser; i = i + 13 {
				for k := i; k < i+13; k++ {
					users[k].DriverUUid = drivers[j].UUid
					UpdateDriverType(drivers[j].UUid, constant.DRIVER_WORKING)
				}
				j++
				if j == len(drivers) {
					break
				}
			}
			if len(drivers) > j {
				//此时还有剩余车辆，用来装剩余的人
				for i := numOfUsers - 1; i >= numOfUsers-numberOfLeftUser; i-- {
					users[i].DriverUUid = drivers[j+1].UUid
					UpdateDriverType(drivers[j+1].UUid, constant.DRIVER_WORKING)
				}
			}
		}
		return users
	}
}

/*
	此函数为分配大巴车Bus的算法，其中大巴车的载客数为14-40人，思想为：
		1.如果用户数量少于一辆大巴车Bus的最大载客量40人，那么直接分配一辆就可以
		2.如果用户数量多于一辆大巴车Bus的最大载客量40人，那么求余数，如果余数为0，说明可以按照最大载客量全部坐满，
		  如果余数不为0，那么将最大载客量的整数倍装满后，剩下的余数个再去装
*/
func assignmentBus(users []module.SmartMeeting, drivers []module.Driver) []module.SmartMeeting {
	if len(drivers) < 1 {
		return users
	}
	numOfUsers := len(users)
	//此时待安排的用户数少于40人，因此一辆大巴车便可以装下所有人
	if numOfUsers <= 40 {
		for i := 0; i < numOfUsers; i++ {
			users[i].DriverUUid = drivers[0].UUid
			UpdateDriverType(drivers[0].UUid, constant.DRIVER_WORKING)
		}
		return users
	} else {
		//此时待安排的用户数多于40，需要多辆大巴车
		//余数为0，正好全部坐满
		if numOfUsers%40 == 0 {
			j := 0
			for i := 0; i < numOfUsers-39; i = i + 40 {
				for k := i; k < 40; k++ {
					users[k].DriverUUid = drivers[j].UUid
					UpdateDriverType(drivers[j].UUid, constant.DRIVER_WORKING)
				}
				j++
				if j == len(drivers) {
					break
				}
			}
		} else {
			//余数为其他数
			numberOfLeftUser := numOfUsers % 40
			j := 0
			for i := 0; i <= numOfUsers-numberOfLeftUser; i = i + 40 {
				for k := i; k < i+40; k++ {
					users[k].DriverUUid = drivers[j].UUid
					UpdateDriverType(drivers[j].UUid, constant.DRIVER_WORKING)
				}
				j++
				if j == len(drivers) {
					break
				}
			}
			if len(drivers) > j {
				//此时还有剩余车辆，用来装剩余的人
				for i := numOfUsers - 1; i >= numOfUsers-numberOfLeftUser; i-- {
					users[i].DriverUUid = drivers[j+1].UUid
					UpdateDriverType(drivers[j+1].UUid, constant.DRIVER_WORKING)
				}
			}
		}
		return users
	}
}

//从meeting表中获取数据拼接成map结构
func getDateMap(meetings *[]module.Meeting) *map[string][]module.MeetingDateInfo {
	meetingInfo := make(map[string][]module.MeetingDateInfo)
	for _, meeting := range *meetings {
		//在此处将meeting拆分，拆成一个为去程，一个为回程
		meetingToFrom := partitionMeeting(&meeting)
		for i := 0; i < 2; i++ {
			if i == 0 {
				//此时为去程信息
				if meetingInfo[meetingToFrom[i].StartDate] == nil {
					list := make([]module.MeetingDateInfo, 0)
					tempMeeting := module.MeetingDateInfo{}
					tempMeeting.MeetingInfo = meeting
					tempMeeting.IfPick = 1
					list = append(list, tempMeeting)
					meetingInfo[meeting.StartDate] = list
					continue
				}
				if meetingInfo[meetingToFrom[i].StartDate] != nil {
					tempMeeting := module.MeetingDateInfo{}
					tempMeeting.MeetingInfo = meeting
					tempMeeting.IfPick = 1
					meetingInfo[meeting.StartDate] = append(meetingInfo[meeting.StartDate], tempMeeting)
				}
			}
			if i == 1 {
				//此时为返程
				if meetingInfo[meetingToFrom[i].ReturnDate] == nil {
					list := make([]module.MeetingDateInfo, 0)
					tempMeeting := module.MeetingDateInfo{}
					tempMeeting.MeetingInfo = meeting
					tempMeeting.IfPick = 0
					list = append(list, tempMeeting)
					meetingInfo[meeting.ReturnDate] = list
					continue
				}
				if meetingInfo[meetingToFrom[i].ReturnDate] != nil {
					tempMeeting := module.MeetingDateInfo{}
					tempMeeting.MeetingInfo = meeting
					tempMeeting.IfPick = 0
					meetingInfo[meeting.ReturnDate] = append(meetingInfo[meeting.ReturnDate], tempMeeting)
				}
			}
		}
	}
	return &meetingInfo
}

//将meeting拆分，拆成一个为去程，一个为回程
func partitionMeeting(meeting *module.Meeting) []module.Meeting {
	var result []module.Meeting
	//此处为生成对应的去程信息
	meetingTo := module.Meeting{}
	meetingTo = *meeting
	meetingTo.ReturnDate = ""
	//此处生成为对应的回城信息Start
	meetingFrom := module.Meeting{}
	meetingFrom = *meeting
	meetingFrom.StartDate = ""
	result = append(result, meetingTo)
	result = append(result, meetingFrom)
	fmt.Println(result)
	return result
}

//封装执行更新航班信息函数
func UpdateShift() {
	UpdateShiftInfo("2021-04-17")
}

// UpdateShiftInfo 更新用户航班信息
func UpdateShiftInfo(tableName string) {
	//此时tableName为航班对应的日期值
	db := sql.GetDB()
	var pickUsers []module.SmartMeeting
	var sentUsers []module.SmartMeeting
	db.Table(tableName).Order("pick_time").Where("pick_or_sent = ?", 1).Find(&pickUsers)
	db.Table(tableName).Order("sent_time").Where("pick_or_sent = ?", 0).Find(&sentUsers)
	fmt.Println("此时获取的pickUser: ", pickUsers)
	fmt.Println("此时获取的sentUser: ", sentUsers)
	go updatePickShiftTime(pickUsers, tableName)
	go updateSentShiftTime(sentUsers, tableName)
}

//比较DB中的时间是否和gRPC调用获取的接口时间一致
func compareShiftTime(originalTime, tableName, Shift string, pickOrSent uint32) (error, bool, string) {
	//调用gRPC获取航班信息
	takeOffTime, landingTime := GetShiftInfo(Shift, tableName)
	if takeOffTime == "" || landingTime == "" {
		logger.ZapLogger.Sugar().Errorf("Table: %+v Shift: %+v do grpc failed. ", tableName, Shift)
		return errorInfo.GRPCServiceError, false, ""
	}
	if pickOrSent == 1 {
		//此时为接站
		if landingTime == originalTime {
			return nil, true, ""
		} else {
			return nil, false, landingTime
		}
	} else {
		//此时为送站
		if takeOffTime == originalTime {
			return nil, true, ""
		} else {
			return nil, false, takeOffTime
		}
	}
}

// goroutine协程处理时间更新
func updatePickShiftTime(pickTime []module.SmartMeeting, tableName string) {
	if len(pickTime) == 0 {
		return
	}
	for _, pick := range pickTime {
		err, result, realTime := compareShiftTime(pick.PickTime, tableName, pick.Shift, 1)
		if err != nil || result == true {
			continue
		} else {
			db := sql.GetDB()
			db.Table("orders").Where("start_date = ? AND start_time = ?", tableName, pick.PickTime).Update("start_time", realTime)
			db.Table(tableName).Where("u_uid = ?", pick.UUid).Update("pick_time", realTime)
		}
	}
}

func updateSentShiftTime(sentTime []module.SmartMeeting, tableName string) {
	if len(sentTime) == 0 {
		return
	}
	for _, sent := range sentTime {
		err, result, realTime := compareShiftTime(sent.PickTime, tableName, sent.Shift, 0)
		if err != nil || result == true {
			continue
		} else {
			db := sql.GetDB()
			db.Table("orders").Where("start_date = ? AND start_time = ?", tableName, sent.SentTime).Update("start_time", realTime)
			db.Table(tableName).Where("u_uid = ?", sent.UUid).Update("sent_time", realTime)
		}
	}
}
