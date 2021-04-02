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
	ID        int    //主键id
	UUid      string //司机对应的UUid
	UserUUid  string //对应用户信息表的uuid
	CarNumber string //车牌号
	CarType   int    //车辆类型，使用int类型进行区分
}

// SmartMeeting 根据时间划分的会议信息的DB-Module
type SmartMeeting struct {
	ID              int    //主键id
	UUid            string //对应的UUid
	MeetingUUid     string //对应的会议信息
	UserName        string //用户姓名
	UserPhoneNumber string //用户电话号码
	FromAddress     string //用户的出发地点
	ToAddress       string //用户到达的地方地点
	PickTime        string //用户的出发时间
	SentTime        string //用户的到达时间
	Shift           string //用户入/离的航班信息
	PickOrSent      uint32 //用户是接站还是送站 1为接站，0为送站
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
	ID                int
	MeetingUUid       string `gorm:"column:meeting_uuid"`        //会议的UUid
	MeetingName       string `gorm:"column:meeting_name"`        //会议的名称
	Name              string `gorm:"column:name"`                //用户的姓名
	Level             uint32 `gorm:"column:level"`               //用户级别
	Company           string `gorm:"column:company"`             //用户所属的机构
	Sex               string `gorm:"column:sex"`                 //用户性别
	IdCard            string `gorm:"column:id_card"`             //身份证号码
	PhoneNumber       string `gorm:"column:phone_number"`        //手机号码
	IfOrderHotel      uint32 `gorm:"column:if_order_hotel"`      //是否订购酒店
	IfOrderPlane      uint32 `gorm:"column:if_order_plane"`      //是否订购机票
	StartDate         string `gorm:"column:start_date"`          //用户的出发日期
	StartTime         string `gorm:"column:start_time"`          //用户的出发时间
	StartBeginAddress string `gorm:"column:start_begin_address"` //用户的出发地点
	StartEndAddress   string `gorm:"column:start_end_address"`   //用户到达的地方地点
	StartShift        string `gorm:"column:start_shift"`         //去程航班信息
	ReturnDate        string `gorm:"column:return_date"`         //用户返程的日期
	ReturnTime        string `gorm:"column:return_time"`         //用户的返程出发时间
	ReturnEndAddress  string `gorm:"column:return_end_address"`  //用户回程到达的地方地点
	ReturnShift       string `gorm:"column:return_shift"`        //回程航班信息
	IfSolve           uint32 `gorm:"column:if_solve"`            //是否经过定时任务的处理
}

//从DB中获取信息转化为Map结构时的辅助数据结构，其中，IfPick为0时为送站，IfPick为1时为接站
type MeetingDateInfo struct {
	MeetingInfo Meeting
	IfPick      uint32
}
