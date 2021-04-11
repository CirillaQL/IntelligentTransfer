package service

import (
	"IntelligentTransfer/module"
	"IntelligentTransfer/pkg/mysql"
)

// GeneOrder 根据SmartMeeting-日期表生成对应的订单信息
func GeneOrder(tableName string) {
	//1.找出所有ifOrder为0且Driver不为空的列
	var users []module.SmartMeeting
	db := mysql.GetDB()
	db.Table(tableName).Where("if_order = ? AND driver_u_uid <> ?", 0, "").Find(&users)

}
