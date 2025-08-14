package dao

import (
	"LibraryManagement/dto"
	"LibraryManagement/entity"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserDAO interface {
	CreateUserDAO(db *gorm.DB, user *dto.RegisterDTO) error
	GetUserByUsernameDAO(db *gorm.DB, username string) (*entity.User, error)
	GetUserByIdDAO(db *gorm.DB, id uint) (*entity.User, error)
}

type userDAOImpl struct{}

// UserDao 全局实例
var UserDao UserDAO = &userDAOImpl{}

// CreateUserDAO 创建用户（自动哈希密码）
func (d *userDAOImpl) CreateUserDAO(db *gorm.DB, user *dto.RegisterDTO) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	return db.Create(user).Error
}

// GetUserByUsernameDAO  根据用户名查找用户
func (d *userDAOImpl) GetUserByUsernameDAO(db *gorm.DB, username string) (*entity.User, error) {
	var user entity.User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByIdDAO 根据 ID 查找用户
func (d *userDAOImpl) GetUserByIdDAO(db *gorm.DB, id uint) (*entity.User, error) {
	var user entity.User
	err := db.First(&user, id).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}
