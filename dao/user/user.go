package user

import (
	"GoNexus/common/mysql"
	"GoNexus/model"
	"GoNexus/utils"
	"errors"
	"log"

	"gorm.io/gorm"
)

const (
	UsernameCondition string = "username"
	EmailCondition    string = "email"
)

// IsExistUser 判断用户是否已存在
func IsExistUser(value, condition string) (bool, *model.User) {
	var user *model.User
	var err error
	switch condition {
	case UsernameCondition:
		user, err = mysql.GetUserByUsername(value)
		if errors.Is(err, gorm.ErrRecordNotFound) || user == nil {
			return false, nil
		}
	case EmailCondition:
		user, err = mysql.GetUserByEmail(value)
		if errors.Is(err, gorm.ErrRecordNotFound) || user == nil {
			return false, nil
		}
	}
	return true, user
}

// Register 将用户信息写入数据库
func Register(username, password, email string) (*model.User, error) {
	if user, err := mysql.InsertUser(&model.User{
		Username: username,
		Password: utils.MD5(password),
		Email:    email},
	); err != nil {
		log.Fatalf("Insert user failed: %v", err)
		return nil, err
	} else {
		return user, nil
	}
}
