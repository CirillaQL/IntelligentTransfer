package service

import (
	errorInfo "IntelligentTransfer/error"
	"IntelligentTransfer/module"
	"IntelligentTransfer/pkg/logger"
	sql "IntelligentTransfer/pkg/mysql"
)

// GeneOrder 根据SmartMeeting-日期表生成对应的订单信息
func GeneOrder(tableName string) {
	//1.找出所有ifOrder为0且Driver不为空的列
	var users []module.SmartMeeting
	db := sql.GetDB()
	db.Table(tableName).Where("if_order = ? AND driver_u_uid <> ?", 0, "").Find(&users)
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
func CancelUserOrder(tableName, userPhone string, pickOrSent uint32) error {
	//首先根据tableName与userPhone定位到对应的用户信息，使用pickOrSent确定为接站还是送站
	db := sql.GetDB()
	//存在性校验
	if db.Migrator().HasTable(tableName) == false {
		logger.ZapLogger.Sugar().Errorf("User:%+v cancel order Error. Table doesn't exist", userPhone)
		return errorInfo.TableDoesNotExist
	}
	//更新接送的司机信息
	db.Table(tableName).Where("user_phone = ? ")
	return nil
}

// DriverCancelOrder 司机主动取消订单
func DriverCancelOrder(uuid string) {
	//在Order表中找到对应的数据，删除对应的数据
	db := sql.GetDB()
	db.Table("orders").Where("driver_u_uid = ?", uuid).Update("driver_u_uid", "")
	//更新smartmeeting表中，所有对应的司机id
	db.Table(getToday()).Where("driver_u_uid = ?", uuid).Update("driver_u_uid", "")
}
