package dao

import (
	"LibraryManagement/internal/api"
	"LibraryManagement/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateUserDAO 测试创建用户
func TestCreateUserDAO(t *testing.T) {
	dao, err := setupTestDB()
	assert.NoError(t, err)

	req := &api.RegisterReq{
		Username: "alice",
		Password: "password123",
		Role:     "user",
	}

	err = dao.CreateUserDAO(req)
	assert.NoError(t, err)

	// 验证用户是否创建成功
	var user model.User
	result := dao.db.Where("username = ?", "alice").First(&user)
	assert.NoError(t, result.Error)
	assert.Equal(t, "alice", user.Username)
	assert.NotEmpty(t, user.PasswordHash)
	assert.Equal(t, "user", user.Role)

	// 验证密码是否被哈希
	assert.True(t, len(user.PasswordHash) > 0)
	assert.NotEqual(t, "password123", user.PasswordHash)
}

// TestCreateUserDAO_DuplicateUsername 测试重复用户名
func TestCreateUserDAO_DuplicateUsername(t *testing.T) {
	dao, err := setupTestDB()
	assert.NoError(t, err)

	req := &api.RegisterReq{
		Username: "bob",
		Password: "pass123",
		Role:     "admin",
	}

	// 第一次创建成功
	err = dao.CreateUserDAO(req)
	assert.NoError(t, err)

	// 第二次创建应失败（唯一索引冲突）
	err = dao.CreateUserDAO(req)
	assert.Error(t, err)
}

// TestGetUserByUsernameDAO 测试根据用户名查找用户
func TestGetUserByUsernameDAO(t *testing.T) {
	dao, err := setupTestDB()
	assert.NoError(t, err)

	// 先创建用户
	req := &api.RegisterReq{
		Username: "charlie",
		Password: "secret",
		Role:     "user",
	}
	err = dao.CreateUserDAO(req)
	assert.NoError(t, err)

	// 查询用户
	user, err := dao.GetUserByUsernameDAO("charlie")
	assert.NoError(t, err)
	assert.Equal(t, "charlie", user.Username)
	assert.Equal(t, "user", user.Role)
}

// TestGetUserByUsernameDAO_UserNotFound 测试用户不存在
func TestGetUserByUsernameDAO_UserNotFound(t *testing.T) {
	dao, err := setupTestDB()
	assert.NoError(t, err)

	user, err := dao.GetUserByUsernameDAO("notexist")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "record not found") // GORM 默认错误
}

// TestGetUserByIdDAO 测试根据 ID 查找用户
func TestGetUserByIdDAO(t *testing.T) {
	dao, err := setupTestDB()
	assert.NoError(t, err)

	// 创建用户
	req := &api.RegisterReq{
		Username: "david",
		Password: "pass",
		Role:     "member",
	}
	err = dao.CreateUserDAO(req)
	assert.NoError(t, err)

	// 查询刚创建用户的 ID
	var user model.User
	dao.db.Where("username = ?", "david").First(&user)

	// 调用 DAO 方法
	foundUser, err := dao.GetUserByIdDAO(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, "david", foundUser.Username)
}

// TestGetUserByIdDAO_UserNotFound 测试用户 ID 不存在
func TestGetUserByIdDAO_UserNotFound(t *testing.T) {
	dao, err := setupTestDB()
	assert.NoError(t, err)

	user, err := dao.GetUserByIdDAO(999) // 不存在的 ID
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "用户不存在", err.Error())
}
