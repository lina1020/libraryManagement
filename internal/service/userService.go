package service

import (
	"LibraryManagement/internal/api"
	"LibraryManagement/internal/model"
	"LibraryManagement/internal/repo/dao"
	"LibraryManagement/internal/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type UserService interface {
	CreateUser(user *api.RegisterReq) error
	Login(dto *api.LoginReq) (*api.LoginResp, error)
}

type userServiceImpl struct{}

func NewUserService() UserService {
	return &userServiceImpl{}
}

func (u userServiceImpl) CreateUser(user *api.RegisterReq) error {

	usernameDAO, err := dao.ApiDao.GetUserByUsernameDAO(user.Username)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 真正的数据库错误
		return err
	}

	if usernameDAO != nil {
		// 用户已存在
		return ErrUserExists
	}

	err = dao.ApiDao.CreateUserDAO(user)

	return err

}

func (u userServiceImpl) Login(dto *api.LoginReq) (*api.LoginResp, error) {
	user, err := dao.ApiDao.GetUserByUsernameDAO(dto.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	//生成 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &api.LoginResp{
		Token:  token,
		UserID: user.ID,
		Role:   user.Role,
	}, nil

}

// GetUserByID 获取用户信息（用于中间件或后续接口）
func (u userServiceImpl) GetUserByID(id uint) (*model.User, error) {
	return dao.ApiDao.GetUserByIdDAO(id)
}
