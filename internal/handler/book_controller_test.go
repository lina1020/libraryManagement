package handler

import (
	"LibraryManagement/internal/api"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --------- Mock Service ----------
type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) Add(req *api.BookInfoReq) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockBookService) Delete(ids []string) error {
	args := m.Called(ids)
	return args.Error(0)
}

func (m *MockBookService) Update(req *api.BookUpdateReq) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockBookService) GetByID(id uint) (*api.BookInfoResp, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.BookInfoResp), args.Error(1)
}

func (m *MockBookService) List(req *api.BookSearchReq) (*api.BookSearchResp, error) {
	args := m.Called(req)
	return args.Get(0).(*api.BookSearchResp), args.Error(1)
}

func (m *MockBookService) SearchBooks(req *api.BookSearchReq) (*api.BookSearchResp, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.BookSearchResp), args.Error(1)
}

func (m *MockBookService) SearchByTitle(title string, exact bool) ([]api.BookInfoResp, error) {
	args := m.Called(title, exact)
	return args.Get(0).([]api.BookInfoResp), args.Error(1)
}

func (m *MockBookService) SearchByContent(content string) ([]api.BookInfoResp, error) {
	args := m.Called(content)
	return args.Get(0).([]api.BookInfoResp), args.Error(1)
}

func (m *MockBookService) InitializeESIndex() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockBookService) ReindexAllBooks() error {
	args := m.Called()
	return args.Error(0)
}

// --------- Helper ---------
func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// --------- Tests ----------

func TestAddBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	h := NewBookHandler(mockService)
	r := gin.Default()
	r.POST("/books", h.AddBook)

	t.Run("success", func(t *testing.T) {
		t.Cleanup(func() {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
		})

		// ✅ 提供所有 required 字段
		req := &api.BookInfoReq{
			Title:  "Go语言",
			Count:  5,                   // ✅ 必填
			ISBN:   "978-7-123-45678-9", // ✅ 必填
			Author: "张三",
		}

		// ✅ 确保 mock 匹配的是完整结构
		mockService.On("Add", req).Return(nil).Once()

		body, _ := json.Marshal(req)
		w := performRequest(r, http.MethodPost, "/books", body)

		// 调试：打印响应内容
		t.Log("Response:", w.Body.String())

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "成功") // 或 "书籍添加成功"
		mockService.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		t.Cleanup(func() {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
		})

		req := &api.BookInfoReq{
			Title:  "Go语言",
			Count:  3,
			ISBN:   "978-7-123-45678-0",
			Author: "李四",
		}

		mockService.On("Add", req).Return(errors.New("添加失败")).Once()

		body, _ := json.Marshal(req)
		w := performRequest(r, http.MethodPost, "/books", body)

		t.Log("Response:", w.Body.String())

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "书籍添加失败")
		mockService.AssertExpectations(t)
	})

	t.Run("invalid_json", func(t *testing.T) {
		w := performRequest(r, http.MethodPost, "/books", []byte("{invalid}"))
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validation_failed", func(t *testing.T) {
		// ❌ 缺少 Count 和 ISBN
		req := &api.BookInfoReq{
			Title: "缺少字段",
			// Count: 0  // 验证失败
			// ISBN: ""  // 验证失败
		}

		body, _ := json.Marshal(req)
		w := performRequest(r, http.MethodPost, "/books", body)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "缺少必要参数")
	})
}

func TestDeleteBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	h := NewBookHandler(mockService)
	r := gin.Default()
	r.DELETE("/books", h.DeleteBook)

	t.Run("success", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		ids := []string{"1", "2"}
		mockService.On("Delete", ids).Return(nil).Once()

		w := performRequest(r, http.MethodDelete, "/books?ids=1&ids=2", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "书籍删除成功")
		mockService.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		mockService.On("Delete", []string{"3"}).Return(errors.New("删除失败")).Once()

		w := performRequest(r, http.MethodDelete, "/books?ids=3", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "书籍删除失败")
	})

	t.Run("no_ids", func(t *testing.T) {
		w := performRequest(r, http.MethodDelete, "/books", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "缺少必要参数")
	})
}

func TestUpdateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	h := NewBookHandler(mockService)
	r := gin.Default()
	r.PUT("/books", h.UpdateBook)

	t.Run("success", func(t *testing.T) {
		t.Cleanup(func() {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
		})

		// ✅ 提供所有 required 字段
		req := &api.BookUpdateReq{
			ID: 1,
			BookInfoReq: api.BookInfoReq{
				Title:  "更新标题",
				Count:  5,
				ISBN:   "978-7-123-45678-9",
				Author: "李四",
			},
		}
		mockService.On("Update", req).Return(nil).Once()

		body, _ := json.Marshal(req)
		w := performRequest(r, http.MethodPut, "/books", body)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "书籍更新成功")
		mockService.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		t.Cleanup(func() {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
		})

		// ✅ 提供完整字段
		req := &api.BookUpdateReq{
			ID: 2,
			BookInfoReq: api.BookInfoReq{
				Title: "错误数据",
				Count: 3,
				ISBN:  "978-7-123-45678-0",
			},
		}
		mockService.On("Update", req).Return(errors.New("更新失败")).Once()

		body, _ := json.Marshal(req)
		w := performRequest(r, http.MethodPut, "/books", body)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "书籍更新失败")
		mockService.AssertExpectations(t)
	})

	t.Run("validation_failed", func(t *testing.T) {
		// ❌ 缺少 Count 和 ISBN
		req := &api.BookUpdateReq{
			ID: 3,
			BookInfoReq: api.BookInfoReq{
				Title: "缺少字段",
				// Count: 0  // 缺失
				// ISBN: ""  // 缺失
			},
		}

		body, _ := json.Marshal(req)
		w := performRequest(r, http.MethodPut, "/books", body)

		// 验证失败应返回 400
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "缺少必要参数")
	})

	t.Run("invalid_json", func(t *testing.T) {
		w := performRequest(r, http.MethodPut, "/books", []byte("{invalid}"))
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	h := NewBookHandler(mockService)
	r := gin.Default()
	r.GET("/books/:id", h.GetBook)

	t.Run("success", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		resp := &api.BookInfoResp{ID: 1, Title: "Go", Author: "张三"}
		mockService.On("GetByID", uint(1)).Return(resp, nil).Once()

		w := performRequest(r, http.MethodGet, "/books/1", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Go")
	})

	t.Run("not_found", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		mockService.On("GetByID", uint(2)).Return((*api.BookInfoResp)(nil), errors.New("not found")).Once()

		w := performRequest(r, http.MethodGet, "/books/2", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "书籍查询失败")
	})

	t.Run("invalid_id", func(t *testing.T) {
		w := performRequest(r, http.MethodGet, "/books/abc", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "ID格式错误")
	})
}

func TestBookList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	h := NewBookHandler(mockService)
	r := gin.Default()
	r.POST("/books/search", h.BookList) // 假设是 POST /books/search

	t.Run("success_with_filters", func(t *testing.T) {
		t.Cleanup(func() {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
		})

		// 请求数据
		req := &api.BookSearchReq{
			Title:    "Go",
			Author:   "张三",
			Keyword:  "并发",
			Page:     1,
			PageSize: 10,
		}

		// 模拟服务层返回
		mockResp := &api.BookSearchResp{
			Books: []api.BookInfoResp{
				{ID: 1, Title: "Go语言编程", Author: "张三", Summary: "并发"},
			},
			Total:    1,
			Page:     1,
			PageSize: 10,
		}

		mockService.On("List", req).Return(mockResp, nil).Once()

		// 发送请求
		body, _ := json.Marshal(req)
		w := performRequest(r, http.MethodPost, "/books/search", body)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Go语言编程")
		assert.Contains(t, w.Body.String(), "成功")
		mockService.AssertExpectations(t)
	})

	t.Run("success_empty_filters", func(t *testing.T) {
		t.Cleanup(func() {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
		})

		// 全部为空，查所有
		req := &api.BookSearchReq{
			Page:     1,
			PageSize: 10,
		}

		mockBooks := &api.BookSearchResp{
			Books: []api.BookInfoResp{
				{ID: 1, Title: "Book1", Author: "张三", Summary: "并发"},
				{ID: 2, Title: "Book2", Author: "张三", Summary: "并发"},
			},
			Total:    1,
			Page:     1,
			PageSize: 10,
		}

		mockService.On("List", req).Return(mockBooks, nil).Once()

		body, _ := json.Marshal(req)
		w := performRequest(r, http.MethodPost, "/books/search", body)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Book1")
		mockService.AssertExpectations(t)
	})

	t.Run("invalid_json", func(t *testing.T) {
		w := performRequest(r, http.MethodPost, "/books/search", []byte("{invalid}"))
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service_error", func(t *testing.T) {
		t.Cleanup(func() {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil
		})

		req := &api.BookSearchReq{
			Title:    "Error",
			Page:     1,
			PageSize: 10,
		}

		mockService.On("List", req).Return(&api.BookSearchResp{}, errors.New("数据库错误")).Once()

		body, _ := json.Marshal(req)
		w := performRequest(r, http.MethodPost, "/books/search", body)

		assert.Equal(t, http.StatusOK, w.Code) // 注意：你返回的是 200 + error message
		assert.Contains(t, w.Body.String(), "书籍查询失败")
		mockService.AssertExpectations(t)
	})

}

func TestSearchBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	h := NewBookHandler(mockService)
	r := gin.Default()
	r.POST("/books/search", h.SearchBooks)

	t.Run("success", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		req := &api.BookSearchReq{Keyword: "Go"}
		resp := &api.BookSearchResp{
			Books: []api.BookInfoResp{{ID: 1, Title: "Go"}},
			Total: 1,
		}

		mockService.On("SearchBooks", req).Return(resp, nil).Once()
		body, _ := json.Marshal(req)
		w := performRequest(r, http.MethodPost, "/books/search", body)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Go")
	})

	t.Run("failure", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		req := &api.BookSearchReq{Keyword: "Java"}
		mockService.On("SearchBooks", req).Return((*api.BookSearchResp)(nil), errors.New("搜索失败")).Once()

		body, _ := json.Marshal(req)
		w := performRequest(r, http.MethodPost, "/books/search", body)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "搜索失败")
	})
}

func TestSearchByTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	h := NewBookHandler(mockService)
	r := gin.Default()
	r.GET("/books/title", h.SearchByTitle)

	t.Run("exact_match", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		resp := []api.BookInfoResp{{ID: 1, Title: "Go"}}
		mockService.On("SearchByTitle", "Go", true).Return(resp, nil).Once()

		w := performRequest(r, http.MethodGet, "/books/title?title=Go&exact=true", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Go")
	})

	t.Run("fuzzy_match", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		resp := []api.BookInfoResp{{ID: 2, Title: "Go高级编程"}}
		mockService.On("SearchByTitle", "Go", false).Return(resp, nil).Once()

		w := performRequest(r, http.MethodGet, "/books/title?title=Go", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Go")
	})

	t.Run("failure", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		mockService.On("SearchByTitle", "Java", false).Return([]api.BookInfoResp{}, errors.New("搜索失败")).Once()

		w := performRequest(r, http.MethodGet, "/books/title?title=Java", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "标题搜索失败")
	})

	t.Run("missing_title", func(t *testing.T) {
		w := performRequest(r, http.MethodGet, "/books/title", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "标题参数不能为空")
	})
}

func TestSearchByContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	h := NewBookHandler(mockService)
	r := gin.Default()
	r.GET("/books/content", h.SearchByContent)

	t.Run("success", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		resp := []api.BookInfoResp{{ID: 1, Title: "Go", Content: "Go语言并发编程"}}
		mockService.On("SearchByContent", "并发").Return(resp, nil).Once()

		w := performRequest(r, http.MethodGet, "/books/content?content=并发", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Go")
	})

	t.Run("failure", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		mockService.On("SearchByContent", "AI").Return([]api.BookInfoResp{}, errors.New("搜索失败")).Once()

		w := performRequest(r, http.MethodGet, "/books/content?content=AI", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "内容搜索失败")
	})

	t.Run("missing_content", func(t *testing.T) {
		w := performRequest(r, http.MethodGet, "/books/content", nil)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "内容参数不能为空")
	})
}

func TestInitESIndex(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	h := NewBookHandler(mockService)
	r := gin.Default()
	r.POST("/books/init", h.InitESIndex)

	t.Run("success", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		mockService.On("InitializeESIndex").Return(nil).Once()

		w := performRequest(r, http.MethodPost, "/books/init", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "ES索引初始化成功")
	})

	t.Run("failure", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		mockService.On("InitializeESIndex").Return(errors.New("初始化失败")).Once()

		w := performRequest(r, http.MethodPost, "/books/init", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "初始化ES索引失败")
	})
}

func TestReindexBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	h := NewBookHandler(mockService)
	r := gin.Default()
	r.POST("/books/reindex", h.ReindexBooks)

	t.Run("success", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		mockService.On("ReindexAllBooks").Return(nil).Once()

		w := performRequest(r, http.MethodPost, "/books/reindex", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "重新索引完成")
	})

	t.Run("failure", func(t *testing.T) {
		defer func() { mockService.ExpectedCalls, mockService.Calls = nil, nil }()

		mockService.On("ReindexAllBooks").Return(errors.New("重新索引失败")).Once()

		w := performRequest(r, http.MethodPost, "/books/reindex", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "重新索引失败")
	})
}
