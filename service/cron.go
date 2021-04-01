package service

import (
	"IntelligentTransfer/module"
	"IntelligentTransfer/pkg/logger"
	"IntelligentTransfer/pkg/mysql"
	"fmt"
	"github.com/robfig/cron/v3"
)

var c *cron.Cron

func init() {
	c = cron.New()
	logger.ZapLogger.Sugar().Info("Cron init success")
}

func StartCron() {
	logger.ZapLogger.Sugar().Info("Start Cron service success")
	c.AddFunc("@every 10s", Work)
	c.Start()
	select {}
}

func Work() {
	db := mysql.GetDB()
	var result []module.Meeting
	db.Where("if_solve = ?", 0).Find(&result)
	ans := getDateMap(&result)
	for key, _ := range *ans {
		if db.Migrator().HasTable(key) == false {
			fmt.Println("Ok")
			_ = db.Table(key).Migrator().CreateTable(module.Meeting{})
			//for _, v := range value {
			//	test := module.Meeting{}
			//	test = v.MeetingInfo
			//	db.Create(test)
			//}
		}
	}
}

func getDateMap(meetings *[]module.Meeting) *map[string][]module.MeetingDateInfo {
	var meetingInfo map[string][]module.MeetingDateInfo
	meetingInfo = make(map[string][]module.MeetingDateInfo)
	for _, meeting := range *meetings {
		if meetingInfo[meeting.ToDate] == nil {
			list := make([]module.MeetingDateInfo, 0)
			tempMeeting := module.MeetingDateInfo{}
			tempMeeting.MeetingInfo = meeting
			tempMeeting.IfPick = 1
			list = append(list, tempMeeting)
			meetingInfo[meeting.ToDate] = list
		}
		if meetingInfo[meeting.FromDate] == nil {
			list := make([]module.MeetingDateInfo, 0)
			tempMeeting := module.MeetingDateInfo{}
			tempMeeting.MeetingInfo = meeting
			tempMeeting.IfPick = 0
			list = append(list, tempMeeting)
			meetingInfo[meeting.FromDate] = list
		}
		if meetingInfo[meeting.ToDate] != nil {
			tempMeeting := module.MeetingDateInfo{}
			tempMeeting.MeetingInfo = meeting
			tempMeeting.IfPick = 1
			meetingInfo[meeting.ToDate] = append(meetingInfo[meeting.ToDate], tempMeeting)
		}
		if meetingInfo[meeting.FromDate] != nil {
			tempMeeting := module.MeetingDateInfo{}
			tempMeeting.MeetingInfo = meeting
			tempMeeting.IfPick = 0
			meetingInfo[meeting.FromDate] = append(meetingInfo[meeting.FromDate], tempMeeting)
		}
	}
	return &meetingInfo
}
