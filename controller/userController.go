package controller

import (
	"LibraryManagement/config"
	"LibraryManagement/dto"
	"LibraryManagement/result"
	"LibraryManagement/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

// LoginController 登入
func LoginController(c *gin.Context) {
	loginDTO := &dto.LoginDTO{}
	err := c.BindJSON(loginDTO)
	if err != nil {
		result.Failed(c, result.FailedCode, "请求数据格式错误")
		return
	}
	fmt.Println("收到请求---登入: ", loginDTO)
	service.NewUserService(c, config.DB).Login(loginDTO)

}

// RegisterController 注册
func RegisterController(c *gin.Context) {
	registerDTO := &dto.RegisterDTO{}
	err := c.BindJSON(registerDTO)
	if err != nil {
		result.Failed(c, result.FailedCode, "请求数据格式错误")
		return
	}
	fmt.Println("收到请求---注册: ", registerDTO)
	service.NewUserService(c, config.DB).CreateUser(registerDTO)
}
