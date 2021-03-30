/*
 数据库gorm中使用的DB-struct定义
*/
package module

import "time"

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
	ID        int    //主键id
	UUid      string //司机对应的UUid
	UserUUid  string //对应用户信息表的uuid
	CarNumber string //车牌号
	CarType   int    //车辆类型，使用int类型进行区分
}

// Meeting 会议DB-Module
//type Meeting struct {
//	ID              int       //主键id
//	UUid            string    //对应的会议信息
//	MeetingName     string    //会议的名称
//	UserId          string    //参加会议人员的用户id
//	UserName        string    //用户姓名
//	UserPhoneNumber string    //用户电话号码
//	OrganiserId     string    //组织者的用户id
//	OrganiserName   string    //组织者的姓名
//	IfOrganiser     int       //是否为主办人
//	IfParticipant   int       //是否为讲师等重要参会人
//	LeaveTime       time.Time //去程时间
//	LeaveFromCity   string    //去程城市
//	LeaveToCity     string    //到达城市
//	BackTime        time.Time //返程时间
//	BackFromCity    string    //返程城市
//	BackToCity      string    //返程到达城市
//}

// SmartMeeting 根据时间划分的会议信息的DB-Module
type SmartMeeting struct {
	ID              int       //主键id
	UUid            string    //对应的UUid
	MeetingUUid     string    //对应的会议信息
	UserId          string    //用户的uuid
	UserName        string    //用户姓名
	UserPhoneNumber string    //用户电话号码
	FromAddress     string    //用户的出发地点
	ToAddress       string    //用户到达的地方地点
	BeginTime       time.Time //用户的出发时间
	ArriveTime      time.Time //用户的到达时间
	Shift           string    //用户入/离的航班信息
}

// Order 生成对应的订单的信息
type Order struct {
	ID          int     //主键id
	UUid        string  //订单的UUid
	MeetingUUid string  //该用户参加的会议的uuid
	DriverUUid  string  //关联对应的司机信息
	UserId      string  //关联对应乘客的信息
	FromAddress string  //出发地
	ToAddress   string  //目的地
	Price       float64 //车费
}

// MeetingInfo 从Excel表中获得的数据对应的结构体
type Meeting struct {
	ID               int
	MeetingName      string    //会议的名称
	Name             string    //用户的姓名
	Level            uint32    //用户级别
	Company          string    //用户所属的机构
	Sex              string    //用户性别
	IdCard           string    //身份证号码
	PhoneNumber      string    //手机号码
	IfOrderHotel     uint32    //是否订购酒店
	IfOrderPlane     uint32    //是否订购机票
	ToTime           time.Time //用户的出发时间
	ToBeginAddress   string    //用户的出发地点
	ToEndAddress     string    //用户到达的地方地点
	ToShift          string    //去程航班信息
	FromTime         time.Time //用户的回程出发时间
	FromBeginAddress string    //用户的回城出发地点
	FromEndAddress   string    //用户回程到达的地方地点
	FromShift        string    //回程航班信息
}
