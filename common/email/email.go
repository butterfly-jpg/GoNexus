package email

import (
	"GoNexus/config"
	"fmt"

	"gopkg.in/gomail.v2"
)

const (
	SubjectMsg  = "来自GoNexus的信息"
	CaptchaMsg  = "GoNexus的验证码如下(验证码仅限于2分钟有效):"
	UserNameMsg = "GoNexus的账号如下，请保留好，后续可以用账号进行登录"
)

// SendCaptcha 向邮箱发送验证码
func SendCaptcha(email, captcha, msg string) error {
	// 1. 创建新邮件对象
	m := gomail.NewMessage()
	// 2. 设置发件人
	m.SetHeader("From", config.GetConfig().EmailConfig.Email)
	// 3. 设置收件人
	m.SetHeader("To", email)
	// 4. 设置主题
	m.SetHeader("Subject", SubjectMsg)
	// 5. 设置正文
	m.SetBody("text/plain", msg+" "+captcha) //纯文本邮件,加载快
	// 6. 配置SMTP服务器消息：QQ邮箱
	dialer := gomail.NewDialer("smtp.qq.com", 587, config.GetConfig().EmailConfig.Email, config.GetConfig().EmailConfig.Authcode)
	// 7. 发送邮件
	if err := dialer.DialAndSend(m); err != nil {
		fmt.Printf("DialAndSend failed. err: %v", err)
		return err
	}
	return nil
}
