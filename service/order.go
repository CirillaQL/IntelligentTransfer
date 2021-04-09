package service

import (
	"IntelligentTransfer/module"
	"IntelligentTransfer/pkg/mysql"
	MyTimeParse "IntelligentTransfer/pkg/time"
)

// GenerateOrders 根据传入的司机信息和乘客表生成订单
func GenerateOrders(driver module.Driver, Passengers []module.Passenger) []module.Order {
	orderList := make([]module.Order, 0)
	for _, passenger := range Passengers {
		order := module.Order{}
		order.UUid = generateUUID()
		order.DriverUUid = driver.UUid
		order.UserName = passenger.UserName
		order.UserPhone = passenger.UserPhone
		order.UserShift = passenger.UserShift
		if order.StartTime == "" {
			order.StartTime = passenger.ArriveTime
		} else {
			olderTime := MyTimeParse.TimeCompareLater(order.StartTime, passenger.ArriveTime)
			order.StartTime = olderTime
		}
		order.CarNumber = driver.CarNumber
		order.CarType = driver.CarType
		order.Price = 0

		orderList = append(orderList, order)
	}
	return orderList
}

// UserToPassenger 用户转为乘客
func UserToPassenger(users []module.SmartMeeting) []module.Passenger {
	//将用户转化为乘客信息
	var passengers []module.Passenger
	for _, user := range users {
		passenger := module.Passenger{}
		passenger.UserName = user.UserName
		passenger.UserPhone = user.UserPhoneNumber
		passenger.UserShift = user.Shift
		if user.PickOrSent == 1 {
			passenger.ArriveTime = user.PickTime
		} else {
			passenger.ArriveTime = user.SentTime
		}
		passengers = append(passengers, passenger)
	}
	return passengers
}

/*
	GenerateMasterOrder 生成主办人的接送订单; 其中主办人的权重最高，优先安排4座小车，为专车接送
*/
func GenerateMasterOrder(users []module.SmartMeeting) []module.Order {
	//1.找出目前所有小轿车司机
	var drivers []module.Driver
	drivers = GetAllTypeOneDriver()
	if len(drivers) != 0 {
		//存在当前空闲的司机，则安排对应的订单
		return GenerateOrders(drivers[0], UserToPassenger(users))
	} else {
		return nil
	}
}

/*
	SaveOrdersDB 保存订单到本地
*/
func SaveOrdersDB(orders []module.Order) error {
	db := mysql.GetDB()
	db.Create(orders)
	return nil
}
