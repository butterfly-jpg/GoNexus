// Package code 响应状态码
package code

type Code int64

const (
	SuccessCode Code = 1000

	InvalidParamsCode    Code = 2001
	UserExistCode        Code = 2002
	UserNotExistCode     Code = 2003
	InvalidPasswordCode  Code = 2004
	NotMatchPasswordCode Code = 2005
	InvalidTokenCode     Code = 2006
	NotLoginCode         Code = 2007
	InvalidCaptchaCode   Code = 2008
	RecordNotFoundCode   Code = 2009
	IllegalPasswordCode  Code = 2010

	ForbiddenCode Code = 3001

	ServerBusyCode Code = 4001

	AIModelNotFind    Code = 5001
	AIModelCannotOpen Code = 5002
	AIModelFail       Code = 5003

	TTSFail Code = 6001
)

var msg = map[Code]string{
	SuccessCode: "success",

	InvalidParamsCode:    "请求参数错误",
	UserExistCode:        "用户名已存在",
	UserNotExistCode:     "用户不存在",
	InvalidPasswordCode:  "用户名或密码错误",
	NotMatchPasswordCode: "两次密码不一致",
	InvalidTokenCode:     "无效的Token",
	NotLoginCode:         "用户未登录",
	InvalidCaptchaCode:   "验证码错误",
	RecordNotFoundCode:   "记录不存在",
	IllegalPasswordCode:  "密码不合法",

	ForbiddenCode: "权限不足",

	ServerBusyCode: "服务繁忙",

	AIModelNotFind:    "模型不存在",
	AIModelCannotOpen: "无法打开模型",
	AIModelFail:       "模型运行失败",
	TTSFail:           "语音服务失败",
}

func (code Code) Code() int64 {
	return int64(code)
}

func (code Code) Msg() string {
	if m, ok := msg[code]; ok {
		return m
	}
	return msg[ServerBusyCode]
}
