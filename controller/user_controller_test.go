package controller

import (
	"LibraryManagement/dto"
	"LibraryManagement/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserService 模拟服务层
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Login(loginDTO *dto.LoginDTO) {
	m.Called(loginDTO)
}

func (m *MockUserService) CreateUser(registerDTO *dto.RegisterDTO) {
	m.Called(registerDTO)
}

// 用于替换 service.NewUserService
var originalNewUserService = service.NewUserService

func setupTest() (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return w, c
}

// TestLoginController_ValidJSON 测试登录成功
func TestLoginController_ValidJSON(t *testing.T) {
	w, c := setupTest()

	// 构造请求体
	loginDTO := &dto.LoginDTO{
		Username: "alice",
		Password: "123456",
	}
	body, _ := json.Marshal(loginDTO)

	// 创建请求
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// 模拟 UserService
	mockSvc := new(MockUserService)
	mockSvc.On("Login", loginDTO).Return()

	// 替换 NewUserService
	service.NewUserService = func(c *gin.Context, db *gorm.DB) service.UserService {
		return mockSvc
	}
	defer func() {
		service.NewUserService = originalNewUserService
	}()

	// 调用控制器
	LoginController(c)

	// 验证
	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)

	// 可选：验证响应体（假设 Login 返回 Success）
	// 如果 Login 返回的是 { "code": 0, "message": "成功", ... }
	// 你可以进一步解析 w.Body.Bytes()
}

// TestLoginController_InvalidJSON 测试 JSON 格式错误
func TestLoginController_InvalidJSON(t *testing.T) {
	w, c := setupTest()

	// 无效 JSON（缺少引号）
	invalidJSON := []byte(`{"username": "alice", "password":}`)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// 确保 UserService 不会被调用
	mockSvc := new(MockUserService)
	mockSvc.On("Login", mock.Anything).Maybe() // 不应被调用

	service.NewUserService = func(c *gin.Context, db *gorm.DB) service.UserService {
		return mockSvc
	}
	defer func() {
		service.NewUserService = originalNewUserService
	}()

	// 调用控制器
	LoginController(c)

	// 验证返回了 Failed 响应
	//  Gin 的 c.BindJSON() 在解析失败时：会自动返回 400 Bad Request，这是 Gin 的默认行为。
	// 即使你在 if err != nil 中调用了 result.Failed(c, ...)，但可能已经写入了响应头，或者 BindJSON 的错误优先级更高。
	assert.Equal(t, http.StatusBadRequest, w.Code) // 注意：Failed 也是 200
	assert.JSONEq(t, `{
		"code": 501,
		"message": "请求数据格式错误",
		"data": {}
	}`, w.Body.String())

	// 验证 Login 没有被调用
	mockSvc.AssertNotCalled(t, "Login", mock.Anything)
}

// TestRegisterController_ValidJSON
func TestRegisterController_ValidJSON(t *testing.T) {
	w, c := setupTest()

	registerDTO := &dto.RegisterDTO{
		Username: "bob",
		Password: "123456",
	}
	body, _ := json.Marshal(registerDTO)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mockSvc := new(MockUserService)
	mockSvc.On("CreateUser", registerDTO).Return()

	service.NewUserService = func(c *gin.Context, db *gorm.DB) service.UserService {
		return mockSvc
	}
	defer func() {
		service.NewUserService = originalNewUserService
	}()

	RegisterController(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

// TestRegisterController_InvalidJSON
func TestRegisterController_InvalidJSON(t *testing.T) {
	w, c := setupTest()

	invalidJSON := []byte(`{"username": "bob", "password":}`)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mockSvc := new(MockUserService)
	service.NewUserService = func(c *gin.Context, db *gorm.DB) service.UserService {
		return mockSvc
	}
	defer func() {
		service.NewUserService = originalNewUserService
	}()

	RegisterController(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{
		"code": 501,
		"message": "请求数据格式错误",
		"data": {}
	}`, w.Body.String())

	mockSvc.AssertNotCalled(t, "CreateUser", mock.Anything)
}
