package handler

import (
	"LibraryManagement/internal/api"
	"LibraryManagement/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -------- Mock UserService --------
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Login(req *api.LoginReq) (*api.LoginResp, error) {
	args := m.Called(req)
	return args.Get(0).(*api.LoginResp), args.Error(1)
}
func (m *MockUserService) CreateUser(req *api.RegisterReq) error {
	args := m.Called(req)
	return args.Error(0)
}

// -------- Tests --------
func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	r := gin.Default()
	r.POST("/login", handler.Login)

	// 成功
	loginReq := api.LoginReq{Username: "alice", Password: "123456"}
	loginResp := &api.LoginResp{Token: "token123"}
	mockService.On("Login", &loginReq).Return(loginResp, nil).Once()

	body, _ := json.Marshal(loginReq)
	w := performRequest(r, http.MethodPost, "/login", body)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token123")

	// 失败：用户名或密码错误
	mockService.On("Login", &loginReq).Return(&api.LoginResp{}, service.ErrInvalidCredentials).Once()
	body2, _ := json.Marshal(loginReq)
	w2 := performRequest(r, http.MethodPost, "/login", body2)
	assert.Contains(t, w2.Body.String(), "用户名或密码错误")

	// 失败：系统错误
	mockService.On("Login", &loginReq).Return(&api.LoginResp{}, errors.New("db down")).Once()
	body3, _ := json.Marshal(loginReq)
	w3 := performRequest(r, http.MethodPost, "/login", body3)
	assert.Contains(t, w3.Body.String(), "系统错误")
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	r := gin.Default()
	r.POST("/register", handler.Register)

	// 成功：新用户
	registerReq := api.RegisterReq{Username: "bob", Password: "123456", Role: "user"}
	mockService.On("CreateUser", &registerReq).Return(nil).Once()

	body, _ := json.Marshal(registerReq)
	w := performRequest(r, http.MethodPost, "/register", body)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "创建成功")

	// 失败：用户已存在
	mockService.On("CreateUser", &registerReq).Return(service.ErrUserExists).Once()
	body2, _ := json.Marshal(registerReq)
	w2 := performRequest(r, http.MethodPost, "/register", body2)
	assert.Contains(t, w2.Body.String(), "用户名已存在")

	// 失败：其他系统错误
	mockService.On("CreateUser", &registerReq).Return(errors.New("insert fail")).Once()
	body3, _ := json.Marshal(registerReq)
	w3 := performRequest(r, http.MethodPost, "/register", body3)
	assert.Contains(t, w3.Body.String(), "创建失败")
}
