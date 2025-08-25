package dao

import (
	"LibraryManagement/internal/api"
	"LibraryManagement/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// dao/user_dao_test.go
type user struct {
	gorm.Model
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"type:text;default:user;not null"` // 测试专用
}

// setupTestDB 初始化内存数据库并自动迁移表
func setupTestDB() (*dbService, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移模型
	err = db.AutoMigrate(&model.Book{})
	if err != nil {
		return nil, err
	}

	// 自动迁移 User 表
	err = db.AutoMigrate(&user{})
	if err != nil {
		return nil, err
	}

	return &dbService{db: db}, nil
}

func TestBookAddDAO(t *testing.T) {
	dao, err := setupTestDB()
	assert.NoError(t, err)

	req := &api.BookInfoReq{
		Title:   "The Go Programming Language",
		Author:  "Alan A.",
		ISBN:    "978-0134190440",
		Count:   5,
		Summary: "A great book about Go",
		Content: "Detailed content...",
	}

	book, err := dao.BookAddDAO(req)
	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, req.Title, book.Title)
	assert.Equal(t, req.ISBN, book.ISBN)
	assert.Equal(t, uint(1), book.ID) // 自增 ID 从 1 开始
	assert.Equal(t, uint(1), uint(book.Version))
}

func TestBookDeleteDAO(t *testing.T) {
	dao, err := setupTestDB()
	assert.NoError(t, err)

	// 先添加一本书
	_ = dao.db.Create(&model.Book{
		Title: "To Delete",
		ISBN:  "978-1111111111",
		Count: 1,
	})

	// 删除 ID 为 1 的书
	err = dao.BookDeleteDAO([]string{"1"})
	assert.NoError(t, err)

	var count int64
	dao.db.Model(&model.Book{}).Count(&count)
	assert.Equal(t, int64(0), count)

	// 测试删除空列表
	err = dao.BookDeleteDAO([]string{})
	assert.NoError(t, err)
}

func TestBookUpdateDAO_OptimisticLock(t *testing.T) {
	dao, err := setupTestDB()
	assert.NoError(t, err)

	// 插入测试数据
	book := model.Book{
		Title:   "Old Title",
		Author:  "Old Author",
		ISBN:    "978-0000000000",
		Count:   1,
		Version: 1,
	}
	dao.db.Create(&book)

	// 查询确认数据正确插入
	var current model.Book
	dao.db.First(&current, book.ID)
	assert.Equal(t, uint(1), uint(current.Version))

	// 第一次更新
	req := &api.BookUpdateReq{
		ID: book.ID,
		BookInfoReq: api.BookInfoReq{
			Title:   "New Title",
			Author:  "New Author",
			ISBN:    "978-1234567890",
			Count:   10,
			Summary: "Updated summary",
			Content: "Updated content",
		},
	}

	updatedBook, err := dao.BookUpdateDAO(req)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", updatedBook.Title)
	assert.Equal(t, uint(2), uint(updatedBook.Version))

	// TODO 模拟并发更新：使用旧 version (1) 尝试更新
	//oldReq := &api.BookUpdateReq{
	//	ID: book.ID,
	//	BookInfoReq: api.BookInfoReq{
	//		Title: "Conflict Title",
	//	},
	//}
	//
	//_, err = dao.BookUpdateDAO(oldReq)
	//assert.Error(t, err)
	//assert.Contains(t, err.Error(), "更新失败：数据已被其他用户修改")
}

func TestBookListDAO_PaginationAndSearch(t *testing.T) {
	dao, err := setupTestDB()
	assert.NoError(t, err)

	// 插入测试数据
	books := []model.Book{
		{Title: "Go in Action", Author: "William K.", ISBN: "978-1617291410", Count: 3, Summary: "Practical Go"},
		{Title: "Learning Go", Author: "Jon Bodner", ISBN: "978-1098138786", Count: 2, Summary: "Learn Go from scratch"},
		{Title: "The Way to Go", Author: "Ivo Balbaert", ISBN: "978-0986289677", Count: 4, Summary: "Comprehensive Go guide"},
	}
	for i := range books {
		books[i].ID = uint(i + 1)
		dao.db.Create(&books[i])
	}

	t.Run("全量分页查询", func(t *testing.T) {
		req := &api.BookSearchReq{Page: 1, PageSize: 2}
		resp, err := dao.BookListDAO(req)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
		assert.Equal(t, 2, len(resp.Books))
		assert.Equal(t, 2, resp.PageSize)
		assert.Equal(t, 2, resp.TotalPages)
	})

	t.Run("按标题模糊查询", func(t *testing.T) {
		req := &api.BookSearchReq{Title: "Go", Page: 1, PageSize: 10}
		resp, err := dao.BookListDAO(req)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(resp.Books), 2)
		for _, b := range resp.Books {
			assert.Contains(t, b.Title, "Go")
		}
	})

	t.Run("按作者查询", func(t *testing.T) {
		req := &api.BookSearchReq{Author: "Bodner", Page: 1, PageSize: 10}
		resp, err := dao.BookListDAO(req)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Books))
		assert.Equal(t, "Jon Bodner", resp.Books[0].Author)
	})

	t.Run("按ISBN精确查询", func(t *testing.T) {
		req := &api.BookSearchReq{ISBN: "978-1617291410", Page: 1, PageSize: 10}
		resp, err := dao.BookListDAO(req)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Books))
		assert.Equal(t, "978-1617291410", resp.Books[0].ISBN)
	})

	t.Run("无效页码默认处理", func(t *testing.T) {
		req := &api.BookSearchReq{Page: -1, PageSize: -5}
		resp, err := dao.BookListDAO(req)
		assert.NoError(t, err)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 10, resp.PageSize)
	})
}

func TestBookGetByIDDAO(t *testing.T) {
	dao, err := setupTestDB()
	assert.NoError(t, err)

	book := model.Book{Title: "Test Book", ISBN: "978-0000000001", Count: 1}
	dao.db.Create(&book)

	found, err := dao.BookGetByIDDAO(book.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, book.Title, found.Title)
	assert.Equal(t, book.ISBN, found.ISBN)

	_, err = dao.BookGetByIDDAO(999)
	assert.Error(t, err)
}

func TestBookGetByISBNDAO(t *testing.T) {
	dao, err := setupTestDB()
	assert.NoError(t, err)

	isbn := "978-0000000002"
	book := model.Book{Title: "ISBN Book", ISBN: isbn, Count: 1}
	dao.db.Create(&book)

	found, err := dao.BookGetByISBNDAO(isbn)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, isbn, found.ISBN)

	_, err = dao.BookGetByISBNDAO("nonexistent")
	assert.Error(t, err)
}
