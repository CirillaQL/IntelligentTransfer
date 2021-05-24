package service

import (
	"IntelligentTransfer/constant"
	"IntelligentTransfer/module"
	"IntelligentTransfer/pkg/encrypt"
	"IntelligentTransfer/pkg/logger"
	sql "IntelligentTransfer/pkg/mysql"
)

// GetMeetingInfo 根据用户信息和会议日期查看会议信息
func GetMeetingInfo(userId, MeetingDate string) []module.Meeting {
	//根据用户id查询用户信息
	db := sql.GetDB()
	var user module.User
	db.Table("users").Where("uuid = ?", userId).Find(&user)
	user.PhoneNumber, _ = encrypt.AesDecrypt(user.PhoneNumber)
	logger.ZapLogger.Sugar().Info(user.PhoneNumber)
	//在SmartMeeting表中查询用户对应的MeetingID列表
	var smartMeeting []module.SmartMeeting
	db.Table(MeetingDate).Where("user_phone_number = ?", user.PhoneNumber).Scan(&smartMeeting)
	logger.ZapLogger.Sugar().Info(smartMeeting)
	//根据获取到的MeetingUUid获取对应的Meeting信息
	var meetings []module.Meeting
	for _, uuid := range smartMeeting {
		var meeting module.Meeting
		db.Table("meetings").Where("meeting_uuid = ? and phone_number = ?", uuid.MeetingUUid, user.PhoneNumber).Find(&meeting)
		meetings = append(meetings, meeting)
	}
	logger.ZapLogger.Sugar().Info(meetings)
	return meetings
}

// UpdateMeeting 更新会议信息
func UpdateMeeting(userId string, json map[string]interface{}) error {
	logger.ZapLogger.Sugar().Info(json)
	//定位到对应的Meeting信息
	var meeting module.Meeting
	db := sql.GetDB()
	db.Table("meetings").Where("UUid = ?", json["UUid"].(string)).Find(&meeting)
	logger.ZapLogger.Sugar().Infof("Find value: %+v", meeting)
	//删除原Meeting表信息
	db.Table("meetings").Where("UUid = ?", meeting.UUid).Delete(&meeting)
	//删除SmartMeeting表信息
	db.Table(meeting.ReturnDate).Where("meeting_u_uid = ? and user_phone_number = ?", meeting.MeetingUUid,
		meeting.PhoneNumber).Delete(module.SmartMeeting{})
	db.Table(meeting.StartDate).Where("meeting_u_uid = ? and user_phone_number = ?", meeting.MeetingUUid,
		meeting.PhoneNumber).Delete(module.SmartMeeting{})
	//删除Order信息
	var order module.Order
	db.Table("orders").Where("user_phone = ? and user_name = ?", meeting.PhoneNumber,
		meeting.Name).Find(&order)
	db.Table("orders").Where("user_phone = ? and user_name = ?", meeting.PhoneNumber,
		meeting.Name).Delete(&order)
	//更新司机信息
	db.Table("drivers").Where("u_uid = ?", order.DriverUUid).Update("status_now", constant.DRIVER_READY)
	//更新Meeting信息
	meeting.UUid = json["UUid"].(string)
	meeting.Name = json["Name"].(string)
	meeting.PhoneNumber = json["PhoneNumber"].(string)
	meeting.StartDate = json["StartDate"].(string)
	meeting.StartTime = json["StartTime"].(string)
	meeting.StartBeginAddress = json["StartBeginAddress"].(string)
	meeting.StartEndAddress = json["StartEndAddress"].(string)
	meeting.StartShift = json["StartShift"].(string)
	meeting.ReturnDate = json["ReturnDate"].(string)
	meeting.ReturnTime = json["ReturnTime"].(string)
	meeting.ReturnStartAddress = json["ReturnStartAddress"].(string)
	meeting.ReturnEndAddress = json["ReturnEndAddress"].(string)
	meeting.ReturnShift = json["ReturnShift"].(string)
	meeting.IfSolve = 0
	//将更新后的信息保存的DB
	db.Table("meetings").Where("UUid = ?", meeting.UUid).Save(&meeting)
	return nil
}
