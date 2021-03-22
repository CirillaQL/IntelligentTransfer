package service

import (
	"IntelligentTransfer/module"
	"IntelligentTransfer/pkg/encrypt"
	"IntelligentTransfer/pkg/logger"
	"IntelligentTransfer/pkg/mysql"
	uuid "github.com/satori/go.uuid"
)

// generateUUID 生成主键uuid ...
func generateUUID() string {
	return uuid.NewV4().String()
}

// Register 用户注册服务，根据前端传入的json数据进行拼装后保存到数据库
func Register(json map[string]interface{}) (string, error) {
	//1.解析json，拼装成User结构体
	User := assembleUser(json)
	//对User敏感信息进行加密
	var err error
	User.Address, err = encrypt.AesEncrypt(User.Address)
	if err != nil {
		logger.Errorf("User Register failed. User's Address can't encrypt: %+v", err)
		return "", err
	}
	User.PhoneNumber, err = encrypt.AesEncrypt(User.PhoneNumber)
	if err != nil {
		logger.Errorf("User Register failed. User's PhoneNumber can't encrypt: %+v", err)
		return "", err
	}
	User.Company, err = encrypt.AesEncrypt(User.Company)
	if err != nil {
		logger.Errorf("User Register failed. User's Company can't encrypt: %+v", err)
		return "", err
	}
	if User.Email != "" {
		User.Email, err = encrypt.AesEncrypt(User.Email)
		if err != nil {
			logger.Errorf("User Register failed. User's Email can't encrypt: %+v", err)
			return "", err
		}
	}
	if User.Password != "" {
		User.Password, err = encrypt.AesEncrypt(User.Password)
		if err != nil {
			logger.Errorf("User Register failed. User's Password can't encrypt: %+v", err)
			return "", err
		}
	}
	if User.IDCard != "" {
		User.IDCard, err = encrypt.AesEncrypt(User.IDCard)
		if err != nil {
			logger.Errorf("User Register failed. User's IDCard can't encrypt: %+v", err)
			return "", err
		}
	}
	//保存到DB
	db := mysql.GetDB()
	db.Create(&User)
	return User.UUID, nil
}

// LoginWithPassword 用户登录服务，此接口为根据电话/邮箱和密码登录，验证码登录另写
func LoginWithPassword(userInfo, password string, inputType uint32) (bool, error) {
	//inputType为判断输入用户信息为电话或者邮箱，如果为电话，值为1，如果为邮箱，值为2
	db := mysql.GetDB()
	if inputType == 1 {
		//从DB找对应的用户信息
		user := module.User{}
		phoneNumberEncrypt, err := encrypt.AesEncrypt(userInfo)
		if err != nil {
			logger.Errorf("User Login failed. User's PhoneNumber can't encrypt: %+v", err)
			return false, err
		}
		db.Where("phone_number = ?", phoneNumberEncrypt).Find(&user)
		if user.UUID == "" {
			return false, nil
		}
		passwordDecrypt, err := encrypt.AesDecrypt(user.Password)
		if err != nil {
			logger.Errorf("User Login failed. User's Password can't decrypt: %+v", err)
			return false, err
		}
		if password == passwordDecrypt {
			return true, nil
		} else {
			return false, nil
		}
	} else if inputType == 2 {
		user := module.User{}
		emailEncrypt, err := encrypt.AesEncrypt(userInfo)
		if err != nil {
			logger.Errorf("User Login failed. User's Email can't encrypt: %+v", err)
			return false, err
		}
		db.Where("email = ?", emailEncrypt).Find(&user)
		if user.UUID == "" {
			return false, nil
		}
		passwordDecrypt, err := encrypt.AesDecrypt(user.Password)
		if err != nil {
			logger.Errorf("User Login failed. User's Password can't decrypt: %+v", err)
			return false, err
		}
		if password == passwordDecrypt {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}
}

//拼装User结构体
func assembleUser(json map[string]interface{}) *module.User {
	User := &module.User{}
	uuid := generateUUID()
	userName := json["user_name"].(string)
	nickName := json["nick_name"].(string)
	sex := json["sex"].(string)
	province := json["province"].(string)
	city := json["city"].(string)
	address := json["address"].(string)
	company := json["company"].(string)
	phoneNumber := json["phone_number"].(string)
	email := json["email"].(string)
	password := json["password"].(string)
	idCard := json["id_card"].(string)
	User.UUID = uuid
	User.UserName = userName
	User.NickName = nickName
	User.Sex = sex
	User.Province = province
	User.City = city
	User.Address = address
	User.Company = company
	User.PhoneNumber = phoneNumber
	User.Email = email
	User.Password = password
	User.IDCard = idCard
	User.IfVip = 0
	return User
}
