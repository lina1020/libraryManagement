package service

//
//import (
//	"LibraryManagement/internal/api/dto"
//	result2 "LibraryManagement/internal/api/result"
//	"LibraryManagement/internal/model"
//	"LibraryManagement/utils"
//	"encoding/json"
//	"errors"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"github.com/stretchr/testify/require"
//	"golang.org/x/crypto/bcrypt"
//	"gorm.io/gorm"
//)
//
//type MockUserDAO struct {
//	mock.Mock
//}
//
//func (m *MockUserDAO) CreateUserDAO(db *gorm.DB, user *dto.RegisterDTO) error {
//	args := m.Called(db, user)
//	return args.Error(0)
//}
//
//func (m *MockUserDAO) GetUserByUsernameDAO(db *gorm.DB, username string) (*model.User, error) {
//	args := m.Called(db, username)
//	if args.Get(0) == nil {
//		return nil, args.Error(1)
//	}
//	return args.Get(0).(*model.User), args.Error(1)
//}
//
//func (m *MockUserDAO) GetUserByIdDAO(db *gorm.DB, id uint) (*model.User, error) {
//	args := m.Called(db, id)
//	if args.Get(0) == nil {
//		return nil, args.Error(1)
//	}
//	return args.Get(0).(*model.User), args.Error(1)
//}
//
//func TestUserServiceImpl_CreateUser(t *testing.T) {
//	gin.SetMode(gin.TestMode)
//
//	// 定义测试用例
//	tests := []struct {
//		name           string                       // 测试名称
//		setupMock      func(*MockUserDAO, *gorm.DB) // 如何设置 mock
//		inputDTO       *dto.RegisterDTO
//		expectedStatus int
//		expectedMsg    string
//		expectedCode   int
//		shouldPass     bool // 是否预期成功（用于 mock 断言）
//	}{
//		{
//			name: "注册成功",
//			setupMock: func(mockDAO *MockUserDAO, db *gorm.DB) {
//				mockDAO.On("GetUserByUsernameDAO", db, "testuser").
//					Return((*model.User)(nil), gorm.ErrRecordNotFound) // 用户不存在
//				mockDAO.On("CreateUserDAO", db, mock.AnythingOfType("*dto.RegisterDTO")).
//					Return(nil)
//			},
//			inputDTO: &dto.RegisterDTO{
//				Username: "testuser",
//				Password: "123456",
//				Role:     "user",
//			},
//			expectedStatus: http.StatusOK,
//			expectedMsg:    "成功",
//			expectedCode:   result2.SuccessCode,
//			shouldPass:     true,
//		},
//		{
//			name: "用户名已存在",
//			setupMock: func(mock *MockUserDAO, db *gorm.DB) {
//				mock.On("GetUserByUsernameDAO", db, "existuser").
//					Return(&model.User{Username: "existuser"}, nil) // 用户已存在
//				// 不会调用 CreateUserDAO
//			},
//			inputDTO: &dto.RegisterDTO{
//				Username: "existuser",
//				Password: "123456",
//				Role:     "user",
//			},
//			expectedStatus: http.StatusOK,
//			expectedMsg:    "用户已存在",
//			expectedCode:   result2.FailedCode,
//			shouldPass:     true,
//		},
//		{
//			name: "数据库创建失败",
//			setupMock: func(mockDAO *MockUserDAO, db *gorm.DB) {
//				mockDAO.On("GetUserByUsernameDAO", db, "failuser").
//					Return((*model.User)(nil), gorm.ErrRecordNotFound)
//				mockDAO.On("CreateUserDAO", db, mock.AnythingOfType("*dto.RegisterDTO")).
//					Return(errors.New("db error"))
//			},
//			inputDTO: &dto.RegisterDTO{
//				Username: "failuser",
//				Password: "123456",
//				Role:     "user",
//			},
//			expectedStatus: http.StatusOK,
//			expectedMsg:    "user创建失败",
//			expectedCode:   result2.FailedCode,
//			shouldPass:     true,
//		},
//	}
//
//	// 开始循环测试
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			w := httptest.NewRecorder()
//			c, _ := gin.CreateTestContext(w)
//
//			db := &gorm.DB{} // 实际项目中可用 sqlmock
//			mockDAO := new(MockUserDAO)
//
//			// 设置 mock 行为
//			tt.setupMock(mockDAO, db)
//
//			svc := &userServiceImpl{
//				Ctx: c,
//				Db:  db,
//				Dao: mockDAO,
//			}
//
//			// 调用被测方法
//			svc.CreateUser(tt.inputDTO)
//
//			// 验证 HTTP 响应
//			assert.Equal(t, tt.expectedStatus, w.Code)
//
//			var resp result2.Result
//			err := json.Unmarshal(w.Body.Bytes(), &resp)
//			assert.NoError(t, err)
//			assert.Equal(t, tt.expectedMsg, resp.Message)
//			assert.Equal(t, tt.expectedCode, resp.Code)
//
//			// 验证 mock 是否按预期被调用
//			if tt.shouldPass {
//				mockDAO.AssertExpectations(t)
//			}
//
//			// 清理：避免 mock 跨测试污染
//			mockDAO.ExpectedCalls = nil
//			mockDAO.Calls = nil
//		})
//	}
//
//}
//
//// 由于 GenerateToken 是工具函数，我们通过 monkey patch 模拟（或重构为接口）
//// 这里我们采用 monkey patch（使用 github.com/agiledragon/gomonkey/v2 可选）
//// 但为简化，我们先用函数变量替换方式演示
//
//// 在 utils/token.go 中建议改为：
//// var GenerateToken = func(uid uint, role string) (string, error) { ... }
//
//// 假设 utils.GenerateToken 是可替换的变量
//var originalGenerateToken = utils.GenerateToken
//
//func TestUserServiceImpl_Login(t *testing.T) {
//	gin.SetMode(gin.TestMode)
//
//	// 测试用例
//	tests := []struct {
//		name           string
//		setupMock      func(*MockUserDAO, *gorm.DB)
//		setupToken     func() // 模拟 GenerateToken
//		inputDTO       *dto.LoginDTO
//		expectedStatus int
//		expectedMsg    string
//		expectedCode   int
//		expectLoginVO  bool // 是否期望返回 LoginVO（成功时）
//	}{
//		{
//			name: "登录成功",
//			setupMock: func(mock *MockUserDAO, db *gorm.DB) {
//				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
//				mock.On("GetUserByUsernameDAO", db, "alice").
//					Return(&model.User{
//						Username:     "alice",
//						PasswordHash: string(hashedPassword),
//						Role:         "user",
//					}, nil)
//			},
//			setupToken: func() {
//				utils.GenerateToken = func(uid uint, role string) (string, error) {
//					return "mock-jwt-token", nil
//				}
//			},
//			inputDTO: &dto.LoginDTO{
//				Username: "alice",
//				Password: "correct",
//			},
//			expectedStatus: http.StatusOK,
//			expectedMsg:    "成功",
//			expectedCode:   result2.SuccessCode,
//			expectLoginVO:  true,
//		},
//		{
//			name: "用户不存在",
//			setupMock: func(mock *MockUserDAO, db *gorm.DB) {
//				mock.On("GetUserByUsernameDAO", db, "notexist").
//					Return((*model.User)(nil), gorm.ErrRecordNotFound)
//			},
//			setupToken: nil,
//			inputDTO: &dto.LoginDTO{
//				Username: "notexist",
//				Password: "any",
//			},
//			expectedStatus: http.StatusOK, // 注意：返回 200，但业务码是失败
//			expectedMsg:    "用户名或密码错误",
//			expectedCode:   result2.FailedCode,
//			expectLoginVO:  false,
//		},
//		{
//			name: "密码错误",
//			setupMock: func(mock *MockUserDAO, db *gorm.DB) {
//				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
//				mock.On("GetUserByUsernameDAO", db, "alice").
//					Return(&model.User{
//						Username:     "alice",
//						PasswordHash: string(hashedPassword),
//						Role:         "user",
//					}, nil)
//			},
//			setupToken: nil,
//			inputDTO: &dto.LoginDTO{
//				Username: "alice",
//				Password: "wrong",
//			},
//			expectedStatus: http.StatusOK,
//			expectedMsg:    "用户名或密码错误",
//			expectedCode:   result2.FailedCode,
//			expectLoginVO:  false,
//		},
//		{
//			name: "生成 Token 失败",
//			setupMock: func(mock *MockUserDAO, db *gorm.DB) {
//				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
//				mock.On("GetUserByUsernameDAO", db, "alice").
//					Return(&model.User{
//						Username:     "alice",
//						PasswordHash: string(hashedPassword),
//						Role:         "user",
//					}, nil)
//			},
//			setupToken: func() {
//				utils.GenerateToken = func(uid uint, role string) (string, error) {
//					return "", errors.New("token generate failed")
//				}
//			},
//			inputDTO: &dto.LoginDTO{
//				Username: "alice",
//				Password: "correct",
//			},
//			expectedStatus: http.StatusOK,
//			expectedMsg:    "系统错误",
//			expectedCode:   result2.FailedCode,
//			expectLoginVO:  false,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			// Setup
//			w := httptest.NewRecorder()
//			c, _ := gin.CreateTestContext(w)
//
//			db := &gorm.DB{}
//			mockDAO := new(MockUserDAO)
//
//			// 恢复 GenerateToken
//			defer func() {
//				utils.GenerateToken = originalGenerateToken
//			}()
//
//			// 执行 setup
//			if tt.setupToken != nil {
//				tt.setupToken()
//			}
//			tt.setupMock(mockDAO, db)
//
//			svc := &userServiceImpl{
//				Ctx: c,
//				Db:  db,
//				Dao: mockDAO,
//			}
//
//			// 调用被测方法
//			svc.Login(tt.inputDTO)
//
//			// 验证 HTTP 状态码
//			assert.Equal(t, tt.expectedStatus, w.Code)
//
//			// 解析响应
//			var resp result2.Result
//			err := json.Unmarshal(w.Body.Bytes(), &resp)
//			require.NoError(t, err)
//			assert.Equal(t, tt.expectedMsg, resp.Message)
//			assert.Equal(t, tt.expectedCode, resp.Code)
//
//			// 验证是否返回 LoginVO
//			if tt.expectLoginVO {
//				var resp struct {
//					Code    int    `json:"code"`
//					Message string `json:"message"`
//					Data    struct {
//						Token  string `json:"token"`
//						UserID uint   `json:"userID"`
//						Role   string `json:"role"`
//					} `json:"data"`
//				}
//
//				err := json.Unmarshal(w.Body.Bytes(), &resp)
//				require.NoError(t, err)
//
//				// 断言
//				assert.Equal(t, "mock-jwt-token", resp.Data.Token)
//				assert.Equal(t, "user", resp.Data.Role)
//			}
//
//			// 验证 mock 调用
//			mockDAO.AssertExpectations(t)
//		})
//	}
//}
