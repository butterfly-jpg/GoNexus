package user

import (
	"GoNexus/common/code"
	"GoNexus/controller"
	"GoNexus/service/user"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	// LoginRequest 登录请求体
	LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	// LoginResponse 登录响应体
	LoginResponse struct {
		controller.Response
		Token string `json:"token,omitempty"`
	}
	// CaptchaRequest 发送验证码请求体
	CaptchaRequest struct {
		Email string `json:"email" binding:"required"`
	}
	// CaptchaResponse 发送验证码响应体
	CaptchaResponse struct {
		controller.Response
	}
	// RegisterRequest 注册请求体
	RegisterRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Captcha  string `json:"captcha" binding:"required"`
	}
	// RegisterResponse 注册响应体
	RegisterResponse struct {
		controller.Response
		Token string `json:"token,omitempty"`
	}
)

// Register 控制层的用户注册接口
func Register(c *gin.Context) {
	req := &RegisterRequest{}
	res := &RegisterResponse{}
	// 1. 解析前端参数
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.InvalidParamsCode))
		return
	}
	// 2. 调用service层注册账号方法
	token, resCode := user.Register(req.Email, req.Password, req.Captcha)
	if resCode != code.SuccessCode {
		c.JSON(http.StatusOK, res.CodeOf(resCode))
		return
	}
	// 3. 业务执行成功返回给前端
	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res)
}

// HandleCaptcha 发送验证码功能
func HandleCaptcha(c *gin.Context) {
	req := &CaptchaRequest{}
	res := &CaptchaResponse{}
	// 1. 解析参数
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.InvalidParamsCode))
		return
	}
	log.Println(req.Email)
	// 2. 调用方法
	resCode := user.SendCaptcha(req.Email)
	if resCode != code.SuccessCode {
		c.JSON(http.StatusOK, res.CodeOf(resCode))
		return
	}
	// 3. 返回前端
	res.Success()
	c.JSON(http.StatusOK, res)
}

// Login 登录功能
func Login(c *gin.Context) {
	req := &LoginRequest{}
	res := &LoginResponse{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.InvalidParamsCode))
		return
	}
	token, resCode := user.Login(req.Username, req.Password)
	if resCode != code.SuccessCode {
		c.JSON(http.StatusOK, res.CodeOf(resCode))
		return
	}
	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res)
}
