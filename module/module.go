/*
 数据库gorm中使用的DB-struct定义
*/
package module

// User 用户DB-Module
type User struct {
	ID          uint32 //自增主键
	UUID        string `gorm:"column:uuid"`         //用户uuid
	UserName    string `gorm:"column:user_name"`    //用户姓名
	NickName    string `gorm:"column:nickname"`     //用户昵称
	Sex         string `gorm:"column:sex"`          //用户性别
	Province    string `gorm:"column:province"`     //用户所在省份
	City        string `gorm:"column:city"`         //用户所在市
	District    string `gorm:"column:district"`     //用户所在区
	Address     string `gorm:"column:address"`      //用户详细地址
	Company     string `gorm:"column:company"`      //用户公司名称
	PhoneNumber string `gorm:"column:phone_number"` //用户手机号码
	Email       string `gorm:"column:email"`        //用户邮箱
	Password    string `gorm:"column:password"`     //用户密码
	IDCard      string `gorm:"column:ID_card"`      //用户身份证信息
	IfVip       uint32 `gorm:"column:if_vip"`       //用户是否为vip
}

// Driver 司机DB-Module
type Driver struct {
	ID        int     //主键id
	UUid      string  `gorm:"column:u_uid"`      //司机对应的UUid
	UserUUid  string  `gorm:"column:user_u_uid"` //对应用户信息表的uuid
	CarNumber string  `gorm:"column:car_number"` //车牌号
	CarType   float64 `gorm:"column:car_type"`   //车辆类型，使用int类型进行区分
	StatusNow uint32  `gorm:"column:status_now"` //司机当前状态
}

// SmartMeeting 根据时间划分的会议信息的DB-Module
type SmartMeeting struct {
	ID              int    //主键id
	UUid            string //对应的UUid
	MeetingUUid     string //对应的会议信息
	UserName        string //用户姓名
	UserPhoneNumber string //用户电话号码
	LeveL           uint32 //用户的等级
	FromAddress     string //用户的出发地点
	ToAddress       string //用户到达的地方地点
	PickTime        string //用户的出发时间
	SentTime        string //用户的到达时间
	Shift           string //用户入/离的航班信息
	PickOrSent      uint32 //用户是接站还是送站 1为接站，0为送站
	IfOrder         uint32 //是否生成了订单 0为未生成 1为已生成
	DriverUUid      string //司机的uuid
}

// Meeting MeetingInfo 从Excel表中获得的数据对应的结构体
type Meeting struct {
	ID                 int
	UUid               string `gorm:"column:UUid"`                 //用作逻辑处理的UUid
	MeetingUUid        string `gorm:"column:meeting_uuid"`         //会议的UUid
	MeetingName        string `gorm:"column:meeting_name"`         //会议的名称
	Name               string `gorm:"column:name"`                 //用户的姓名
	Level              uint32 `gorm:"column:level"`                //用户级别
	Company            string `gorm:"column:company"`              //用户所属的机构
	Sex                string `gorm:"column:sex"`                  //用户性别
	IdCard             string `gorm:"column:id_card"`              //身份证号码
	PhoneNumber        string `gorm:"column:phone_number"`         //手机号码
	IfOrderHotel       uint32 `gorm:"column:if_order_hotel"`       //是否订购酒店
	IfOrderPlane       uint32 `gorm:"column:if_order_plane"`       //是否订购机票
	StartDate          string `gorm:"column:start_date"`           //用户的出发日期
	StartTime          string `gorm:"column:start_time"`           //用户的出发时间
	StartBeginAddress  string `gorm:"column:start_begin_address"`  //用户的出发地点
	StartEndAddress    string `gorm:"column:start_end_address"`    //用户到达的地方地点
	StartShift         string `gorm:"column:start_shift"`          //去程航班信息
	ReturnDate         string `gorm:"column:return_date"`          //用户返程的日期
	ReturnTime         string `gorm:"column:return_time"`          //用户的返程出发时间
	ReturnStartAddress string `gorm:"column:return_start_address"` //用户回程出发的地方
	ReturnEndAddress   string `gorm:"column:return_end_address"`   //用户回程到达的地方地点
	ReturnShift        string `gorm:"column:return_shift"`         //回程航班信息
	IfSolve            uint32 `gorm:"column:if_solve"`             //是否经过定时任务的处理
}

// MeetingDateInfo 从DB中获取信息转化为Map结构时的辅助数据结构，其中，IfPick为0时为送站，IfPick为1时为接站
type MeetingDateInfo struct {
	MeetingInfo Meeting
	IfPick      uint32
}

// Order 生成对应的订单的信息
type Order struct {
	ID          int     //自增id
	UUid        string  `gorm:"column:uuid"`         //订单的UUid
	DriverUUid  string  `gorm:"column:driver_uuid"`  //关联Driver表的uuid
	UserName    string  `gorm:"column:user_name"`    //乘客的名字
	UserPhone   string  `gorm:"column:user_phone"`   //乘客的电话号码
	UserShift   string  `gorm:"column:user_shift"`   //乘客的航班
	StartDate   string  `gorm:"column:start_date"`   //用户出发的日期
	StartTime   string  `gorm:"column:start_time"`   //用户出发的时间
	FromAddress string  `gorm:"column:from_address"` //用户出发的地点
	ToAddress   string  `gorm:"column:to_address"`   //用户到达的地点
	CarNumber   string  `gorm:"column:car_number"`   //车牌号
	CarType     float64 `gorm:"column:car_type"`     //车辆种类
	Price       float64 `gorm:"column:price"`        //预计车费
	PickOrSent  uint32  `gorm:"column:pick_or_sent"` //接送信息，接站为1，送站为2
}

// Passenger 乘客信息
type Passenger struct {
	UserName   string //用户姓名
	UserPhone  string //用户电话号码
	UserShift  string //用户对应的航班信息
	ArriveTime string //用户到达时间
}

// ShiftInfo 航班信息
type ShiftInfo struct {
	Shift     string //航班号
	StartTime string //起飞时间
	EndTime   string //降落时间
	IfDelay   uint32 //是否延误
}
