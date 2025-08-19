package handler

import (
	"LibraryManagement/internal/api"
	"LibraryManagement/internal/api/result"
	"LibraryManagement/internal/service"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Login 登入
func (u *UserHandler) Login(c *gin.Context) {
	loginReq := &api.LoginReq{}
	err := c.BindJSON(loginReq)
	if err != nil {
		result.Failed(c, result.FailedCode, "请求数据格式错误")
		return
	}
	log.Println("收到请求---登入: ", loginReq)

	// 调用服务层
	loginResp, err := u.userService.Login(loginReq)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			result.Failed(c, result.FailedCode, "用户名或密码错误")
		default:
			result.Failed(c, result.FailedCode, "系统错误")
		}
		return
	}

	result.Success(c, loginResp)
}

// Register 注册
func (u *UserHandler) Register(c *gin.Context) {
	registerReq := &api.RegisterReq{}
	err := c.BindJSON(registerReq)
	if err != nil {
		result.Failed(c, result.FailedCode, "请求数据格式错误")
		return
	}
	log.Println("收到请求---注册: ", registerReq)

	if registerReq.Role == "" {
		registerReq.Role = "user"
	}

	// 验证器会根据结构体里写的 validate 标签，自动检查字段是否符合规则
	err = validator.New().Struct(registerReq)
	if err != nil {
		result.Failed(c, result.RequiredCode, result.GetMessage(result.RequiredCode))
		return
	}

	err = u.userService.CreateUser(registerReq)

	if err != nil {
		// 区分错误类型
		switch {
		case errors.Is(err, service.ErrUserExists):
			result.Failed(c, result.FailedCode, "用户名已存在")
		default:
			result.Failed(c, result.FailedCode, registerReq.Role+"创建失败")
		}
		return
	}

	result.Success(c, registerReq.Role+"创建成功")

}
