package service

import (
	"IntelligentTransfer/constant"
	"IntelligentTransfer/module"
	"IntelligentTransfer/pkg/logger"
	sql "IntelligentTransfer/pkg/mysql"
	"errors"
)

// GeneOrder 根据SmartMeeting-日期表生成对应的订单信息
func GeneOrder(tableName string) {
	//1.找出所有ifOrder为0且Driver不为空的列
	var users []module.SmartMeeting
	db := sql.GetDB()
	db.Table(tableName).Where("if_order = ? AND driver_u_uid <> ?", 0, "").Find(&users)
	if len(users) != 0 {
		logger.ZapLogger.Sugar().Infof("find unOrder user: %+v", users)
	}
	for _, user := range users {
		//每个人都有自己的一份订单
		order := CreateOrder(user, tableName)
		if db.Migrator().HasTable("orders") {
			db.Create(&order)
		} else {
			_ = db.Migrator().CreateTable(&module.Order{})
			db.Create(&order)
		}
		updateUserOrder(tableName, user.UUid)
	}
}

//CreateOrder 根据传入的user(smartMeeting)生成Order结构体
func CreateOrder(user module.SmartMeeting, tableName string) module.Order {
	var order module.Order
	order.UUid = generateUUID()
	order.DriverUUid = user.DriverUUid
	order.UserName = user.UserName
	order.UserPhone = user.UserPhoneNumber
	order.UserShift = user.Shift
	order.StartDate = tableName
	if user.PickOrSent == 1 {
		//接站
		order.StartTime = user.PickTime
	} else {
		//送站
		order.StartTime = user.SentTime
	}
	if user.PickOrSent == 0 {
		//送站
		order.PickOrSent = 0
		order.FromAddress = user.FromAddress
		order.ToAddress = user.ToAddress
	}
	var driver module.Driver
	db := sql.GetDB()
	db.Where("u_uid = ?", user.DriverUUid).Find(&driver)
	order.CarNumber = driver.CarNumber
	order.CarType = driver.CarType
	if order.CarType == 1 {
		order.Price = 150
	} else if order.CarType == 2 {
		order.Price = 100
	} else if order.CarType == 3 {
		order.Price = 60
	} else if order.Price == 4 {
		order.Price = 20
	} else {
		order.Price = 0
	}
	return order
}

//更新用户的if_order
func updateUserOrder(tableName, uuid string) {
	db := sql.GetDB()
	db.Table(tableName).Where("u_uid = ?", uuid).Update("if_order", 1)
}

// CancelUserOrder 用户主动取消订单,通过传入的信息来定位到对应的用户接送站信息
func CancelUserOrder(uuid string) error {
	db := sql.GetDB()
	var order module.Order
	db.Table("orders").Where("uuid = ?", uuid).Find(&order)
	deleteRow := db.Table("orders").Where("uuid = ?", uuid).Delete(&order)
	logger.ZapLogger.Sugar().Infof("delete row: %+v", deleteRow.RowsAffected)
	logger.ZapLogger.Sugar().Infof("Orderinfo:%+v", order)
	//更新SmartMeeting表中信息
	//db.Table(order.StartDate).Where("user_name = ? and user_phone_number = ?", order.UserName, order.UserPhone).Update("if_order", 0)
	var smart module.SmartMeeting
	db.Table(order.StartDate).Where("user_name = ? and user_phone_number = ?", order.UserName, order.UserPhone).Find(&smart)
	logger.ZapLogger.Sugar().Infof("smartMeeting info: %+v", smart)
	result := db.Table(order.StartDate).Where("u_uid = ?", smart.UUid).Update("if_order", 0)
	logger.ZapLogger.Sugar().Info(result.RowsAffected)
	//更新司机信息
	var drivers []module.Order
	db.Table("orders").Where("driver_uuid = ?", order.DriverUUid).Find(&drivers)
	if len(drivers) == 0 {
		//此时没有这个司机的订单
		UpdateDriverType(order.DriverUUid, constant.DRIVER_BUSY)
	}
	return nil
}

// GetOrders 获取该用户的所用订单
func GetOrders(userPhone string) []module.Order {
	//从Order表中获取所有用户id的订单
	db := sql.GetDB()
	var orderList []module.Order
	db.Table("orders").Where("user_phone = ?", userPhone).Find(&orderList)
	return orderList
}

// UpdateOrder 更新用户订单信息
func UpdateOrder(UUid, userName, userPhone string) error {
	//根据UUid查找指定UUid的订单
	db := sql.GetDB()
	var order module.Order
	db.Table("orders").Where("uuid = ?", UUid).Find(&order)
	if order.UUid == "" {
		logger.ZapLogger.Sugar().Errorf("Can't find this Order")
		err := errors.New("Can't find order")
		return err
	}
	db.Table("orders").Where("uuid = ?", UUid).Update("user_name", userName).Update("user_phone", userPhone)
	return nil
}

// GetDriverOrder 获取司机订单信息
func GetDriverOrder(uuid string) []module.Order {
	//首先根据用户的uuid查找司机的uuid
	var driver module.Driver
	db := sql.GetDB()
	db.Table("drivers").Where("user_u_uid = ?", uuid).Find(&driver)
	//此时找到了对应的司机uuid
	var orders []module.Order
	db.Table("orders").Where("driver_uuid = ?", driver.UUid).Find(&orders)
	return orders
}
