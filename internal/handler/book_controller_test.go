package handler

//
//import (
//	result2 "LibraryManagement/internal/api/result"
//	"LibraryManagement/internal/model"
//	"LibraryManagement/internal/repo/dao"
//	"bytes"
//	"encoding/json"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"gorm.io/driver/sqlite"
//	"gorm.io/gorm"
//)
//
//// TODO 好像写成了集成测试，同时测试Controller 层和Service 层
//
///*
//Controller 层测试
//
//HTTP 请求测试：测试各种HTTP请求格式和参数
//路由测试：测试路由配置是否正确
//参数验证测试：测试无效参数的处理
//集成测试：测试完整的请求-响应流程
//并发测试：测试控制器在高并发下的表现
//错误处理测试：测试各种异常情况
//*/
//
//// 或者使用MockBookService 模拟 BookService
//
//func setupTestDB() *gorm.DB {
//	// 使用 SQLite 内存数据库进行测试
//	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
//	if err != nil {
//		panic("failed to connect database")
//	}
//
//	// 迁移 schema
//	err = db.AutoMigrate(&model.Book{})
//	if err != nil {
//		panic("failed to migrate database")
//		return nil
//	}
//
//	return db
//}
//
//// 设置测试环境
//func SetupTestRouter() *gin.Engine {
//	// 1. 初始化测试数据库
//	dao.DB = setupTestDB()
//
//	gin.SetMode(gin.TestMode)
//	r := gin.New()
//
//	// 注册路由
//	r.POST("/books/add", AddBook)
//	r.DELETE("/books/delete", DeleteBook)
//	r.PUT("/books/update", UpdateBook)
//	r.POST("/books/list", BookList)
//
//	r.POST("/auth/login", Login)
//	r.POST("/auth/register", Register)
//	return r
//}
//
//func TestAddBookController(t *testing.T) {
//	router := SetupTestRouter()
//
//	tests := []struct {
//		name             string
//		payload          interface{}
//		expectedRespCode int    // 期望的 Result.Code
//		expectedMessage  string // 期望的 Result.Message
//	}{
//		{
//			name: "成功添加书籍",
//			payload: dto.BookInfoDTO{
//				Title: "Go语言编程",
//				Count: 10,
//				ISBN:  "9787111558422",
//			},
//			expectedRespCode: result2.SuccessCode,
//			expectedMessage:  "成功",
//		},
//		{
//			name: "缺少必填字段",
//			payload: dto.BookInfoDTO{
//				Title: "",
//				Count: 10,
//				ISBN:  "9787111558422",
//			},
//			expectedRespCode: result2.RequiredCode,
//			expectedMessage:  "缺少必要参数",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			var jsonPayload []byte
//			var err error
//
//			if str, ok := tt.payload.(string); ok {
//				jsonPayload = []byte(str)
//			} else {
//				jsonPayload, err = json.Marshal(tt.payload)
//				assert.NoError(t, err)
//			}
//
//			req, _ := http.NewRequest("POST", "/books/add", bytes.NewBuffer(jsonPayload))
//			req.Header.Set("Content-Type", "application/json")
//
//			w := httptest.NewRecorder()
//			router.ServeHTTP(w, req)
//
//			// 断言 HTTP 状态码始终为 200
//			assert.Equal(t, http.StatusOK, w.Code)
//
//			// 解析响应体
//			var res result2.Result
//			err = json.Unmarshal(w.Body.Bytes(), &res)
//			assert.NoError(t, err)
//
//			// 断言自定义状态码和消息
//			assert.Equal(t, tt.expectedRespCode, res.Code)
//			assert.Equal(t, tt.expectedMessage, res.Message)
//		})
//	}
//}
//
//func TestDeleteBookController(t *testing.T) {
//	router := SetupTestRouter()
//
//	tests := []struct {
//		name             string
//		queryParams      string // 查询参数格式: ids=1&ids=2
//		expectedRespCode int
//		expectedMessage  string
//	}{
//		{
//			name:             "成功删除书籍",
//			queryParams:      "ids=1&ids=2",
//			expectedRespCode: result2.SuccessCode,
//			expectedMessage:  "成功",
//		},
//		{
//			name:             "未提供ids",
//			queryParams:      "",
//			expectedRespCode: result2.RequiredCode,
//			expectedMessage:  "缺少必要参数",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			url := "/books/delete"
//			if tt.queryParams != "" {
//				url += "?" + tt.queryParams
//			}
//
//			req, _ := http.NewRequest("DELETE", url, nil)
//			w := httptest.NewRecorder()
//			router.ServeHTTP(w, req)
//
//			assert.Equal(t, http.StatusOK, w.Code)
//
//			var res result2.Result
//			err := json.Unmarshal(w.Body.Bytes(), &res)
//			assert.NoError(t, err)
//
//			assert.Equal(t, tt.expectedRespCode, res.Code)
//			assert.Equal(t, tt.expectedMessage, res.Message)
//		})
//	}
//}
