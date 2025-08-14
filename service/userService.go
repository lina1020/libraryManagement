package service

import (
	"LibraryManagement/dao"
	"LibraryManagement/dto"
	"LibraryManagement/entity"
	"LibraryManagement/result"
	"LibraryManagement/utils"
	"LibraryManagement/vo"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(user *dto.RegisterDTO)
	Login(dto *dto.LoginDTO)
}

type userServiceImpl struct {
	Ctx *gin.Context
	Db  *gorm.DB
	Dao dao.UserDAO // 依赖注入
}

func (u userServiceImpl) CreateUser(user *dto.RegisterDTO) {
	if user.Role == "" {
		user.Role = "user"
	}

	// 验证器会根据结构体里写的 validate 标签，自动检查字段是否符合规则
	err := validator.New().Struct(user)
	if err != nil {
		result.Failed(u.Ctx, result.RequiredCode, result.GetMessage(result.RequiredCode))
		return
	}

	usernameDAO, err := u.Dao.GetUserByUsernameDAO(u.Db, user.Username)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 真正的数据库错误
		result.Failed(u.Ctx, result.FailedCode, "创建用户失败")
		return
	}

	if usernameDAO != nil {
		// 用户已存在
		result.Failed(u.Ctx, result.FailedCode, "用户已存在")
		return
	}

	err = u.Dao.CreateUserDAO(u.Db, user)
	if err != nil {
		result.Failed(u.Ctx, result.FailedCode, user.Role+"创建失败")
		return
	}

	result.Success(u.Ctx, user.Role+"创建成功")
}

func (u userServiceImpl) Login(dto *dto.LoginDTO) {
	user, err := u.Dao.GetUserByUsernameDAO(u.Db, dto.Username)
	if err != nil {
		result.Failed(u.Ctx, result.FailedCode, "用户名或密码错误")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password)); err != nil {
		result.Failed(u.Ctx, result.FailedCode, "用户名或密码错误")
		return
	}

	//生成 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		result.Failed(u.Ctx, result.FailedCode, "系统错误")
		return
	}

	result.Success(u.Ctx, &vo.LoginVO{
		Token:  token,
		UserID: user.ID,
		Role:   user.Role,
	})

}

// GetUserByID 获取用户信息（用于中间件或后续接口）
func (u userServiceImpl) GetUserByID(id uint) (*entity.User, error) {
	return dao.UserDao.GetUserByIdDAO(u.Db, id)
}

var NewUserService = func(ctx *gin.Context, db *gorm.DB) UserService {
	return &userServiceImpl{
		Ctx: ctx,
		Db:  db,
		Dao: dao.UserDao,
	}
}
