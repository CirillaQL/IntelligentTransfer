/*
 数据库gorm中使用的DB-struct定义
*/
package module

type User struct {
	ID          uint32 //自增主键
	UUID        string `gorm:"column:uuid"`         //用户uuid
	UserName    string `gorm:"column:user_name"`    //用户姓名
	NickName    string `gorm:"column:nickname"`     //用户昵称
	Sex         string `gorm:"column:sex"`          //用户性别
	Province    string `gorm:"column:province"`     //用户所在省份
	City        string `gorm:"column:city"`         //用户所在城市
	Address     string `gorm:"column:address"`      //用户详细地址
	Company     string `gorm:"column:company"`      //用户公司名称
	PhoneNumber string `gorm:"column:phone_number"` //用户手机号码
	Email       string `gorm:"column:email"`        //用户邮箱
	Password    string `gorm:"column:password"`     //用户密码
	IDCard      string `gorm:"column:ID_card"`      //用户身份证信息
	IfVip       uint32 `gorm:"column:if_vip"`       //用户是否为vip
}
