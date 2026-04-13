package user

import (
	"GoNexus/common/code"
	myemail "GoNexus/common/email"
	"GoNexus/common/redis"
	"GoNexus/dao/user"
	"GoNexus/model"
	"GoNexus/utils"
)

// Register 注册账号
func Register(email, password, captcha string) (string, code.Code) {
	// 1. 根据邮箱检验用户是否已存在
	if ok, _ := user.IsExistUser(email, user.EmailCondition); ok {
		return "", code.UserExistCode
	}
	// 2. 读取缓存校验验证码是否有效
	if ok, _ := redis.CheckCaptchaForEmail(email, captcha); ok {
		return "", code.InvalidCaptchaCode
	}
	// 3. 生成11位随机账号
	username := utils.GetRandomNumbers(11)
	// 4. 用户注册信息存入数据库
	userInfo, err := user.Register(username, password, email)
	if err != nil {
		return "", code.ServerBusyCode
	}
	// 5. 将账号信息再发送到申请注册的邮箱中，方便用户查看账号登录
	if err = myemail.SendCaptcha(email, username, myemail.UserNameMsg); err != nil {
		return "", code.ServerBusyCode
	}
	// 6. 生成JWT token
	token, err := utils.GenerateToken(userInfo.ID, userInfo.Username)
	if err != nil {
		return "", code.ServerBusyCode
	}
	return token, code.SuccessCode
}

// SendCaptcha 向指定邮箱发送验证码
func SendCaptcha(email string) code.Code {
	// 1. 生成随机6位数字为验证码
	captcha := utils.GetRandomNumbers(6)
	// 2. 邮箱为key验证码为value先存入redis做缓存
	if err := redis.SetCaptchaForEmail(email, captcha); err != nil {
		return code.ServerBusyCode
	}
	// 3. 远程发送验证码到指定邮箱
	if err := myemail.SendCaptcha(email, captcha, myemail.CaptchaMsg); err != nil {
		return code.ServerBusyCode
	}
	return code.SuccessCode
}

// Login 登录账号
func Login(username, password string) (string, code.Code) {
	var userInfo *model.User
	var ok bool
	// 1. 检查用户名是否存在,支持用户名和邮箱两种方式登录
	if ok, userInfo = user.IsExistUser(username, user.UsernameCondition); !ok {
		if ok, userInfo = user.IsExistUser(username, user.EmailCondition); !ok {
			return "", code.UserNotExistCode
		}
	}
	// 2. 验证密码
	if userInfo.Password != utils.MD5(password) {
		return "", code.InvalidPasswordCode
	}
	// 3. 生成JWT token
	token, err := utils.GenerateToken(userInfo.ID, userInfo.Username)
	if err != nil {
		return "", code.ServerBusyCode
	}
	return token, code.SuccessCode
}
