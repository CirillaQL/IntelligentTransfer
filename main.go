package main

import "IntelligentTransfer/pkg/mysql"

type driver struct {
	UUid      string `gorm:"column:uuid"`
	UserUUid  string `gorm:"column:user_uuid"`
	CarNumber string `gorm:"column:car_number"`
	CarType   uint32 `gorm:"column:car_type"`
}

func main() {
	db := mysql.GetDB()
	defer db.Close()
	dirv := driver{
		UUid:      "dsadasdadas",
		UserUUid:  "dsadasda",
		CarNumber: "dsadadasd",
		CarType:   1,
	}
	db.Create(&dirv)
}
