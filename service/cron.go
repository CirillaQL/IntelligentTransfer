package service

import (
	"IntelligentTransfer/module"
	"IntelligentTransfer/pkg/logger"
	"IntelligentTransfer/pkg/mysql"
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

//开启定时任务
func StartCron() {
	logger.ZapLogger.Sugar().Info("Start Cron service success")
	c.AddFunc("@every 10s", CreateTable)
	c.AddFunc("@every 5s", GetTodayPickInfoByOrder)
	c.Start()
	select {}
}

//生成以日期为TableName的数据库表
func CreateTable() {
	db := mysql.GetDB()
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
					smartmeeting.FromAddress = ""
					smartmeeting.ToAddress = meeting.MeetingInfo.ReturnEndAddress
					smartmeeting.SentTime = meeting.MeetingInfo.ReturnTime
					smartmeeting.Shift = meeting.MeetingInfo.ReturnShift
					smartmeeting.PickOrSent = 0
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

//按照时间从表中获取数据，进行排序，目前获取数据仅仅为获取到当天的数据，其他之后时间段的暂时并不需要
func GetTodayPickInfoByOrder() {
	dateNow := getToday()
	db := mysql.GetDB()
	smartMeeting := make([]module.SmartMeeting, 0)
	if db.Migrator().HasTable(dateNow) {
		//此时获取的为接站
		db.Table(dateNow).Order("sent_time").Where("pick_or_sent = ?", 0).Find(&smartMeeting)
		//获取到按照时间进行了排序的接站信息表
		for _, v := range smartMeeting {
			fmt.Println(v)
		}
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

//分配主办人车辆