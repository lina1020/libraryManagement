package dao

import (
	"LibraryManagement/internal/api"
	"LibraryManagement/internal/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type userDAO interface {
	CreateUserDAO(req *api.RegisterReq) error
	GetUserByUsernameDAO(username string) (*model.User, error)
	GetUserByIdDAO(id uint) (*model.User, error)
}

// CreateUserDAO 创建用户（自动哈希密码）
func (d *dbService) CreateUserDAO(req *api.RegisterReq) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashed)

	user := &model.User{
		Username:     req.Username,
		PasswordHash: string(hashed),
		Role:         req.Role,
	}

	return d.db.Create(user).Error
}

// GetUserByUsernameDAO  根据用户名查找用户
func (d *dbService) GetUserByUsernameDAO(username string) (*model.User, error) {
	var user model.User
	err := d.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByIdDAO 根据 ID 查找用户
func (d *dbService) GetUserByIdDAO(id uint) (*model.User, error) {
	var user model.User
	err := d.db.First(&user, id).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}
