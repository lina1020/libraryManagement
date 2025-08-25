package dao

import (
	"LibraryManagement/internal/api"
	"LibraryManagement/internal/model"
	"errors"
	"fmt"
	"strconv"
)

type bookDAO interface {
	BookAddDAO(req *api.BookInfoReq) (*model.Book, error)
	BookDeleteDAO(idStr []string) error
	BookUpdateDAO(req *api.BookUpdateReq) (*model.Book, error)
	BookListDAO(req *api.BookSearchReq) (*api.BookSearchResp, error)

	BookGetByIDDAO(id uint) (*model.Book, error)
	BookGetByISBNDAO(isbn string) (*model.Book, error)
}

func (d *dbService) BookAddDAO(req *api.BookInfoReq) (*model.Book, error) {
	book := &model.Book{
		Title: req.Title,
		Count: req.Count,
		ISBN:  req.ISBN,

		Author:  req.Author,
		Content: req.Content,
		Summary: req.Summary,
	}

	err := d.db.Create(book).Error
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (d *dbService) BookDeleteDAO(idStr []string) error {
	if len(idStr) == 0 {
		return nil // 没有 ID，无需删除
	}

	ids := make([]uint, 0, len(idStr))
	for _, id := range idStr {
		uid, _ := strconv.ParseUint(id, 10, 64)
		ids = append(ids, uint(uid))
	}

	if len(ids) == 0 {
		return nil // 所有 ID 都非法，不执行删除
	}
	return d.db.Delete(&model.Book{}, ids).Error
}

// 数据库乐观锁更新
func (d *dbService) BookUpdateDAO(req *api.BookUpdateReq) (*model.Book, error) {
	var book model.Book
	err := d.db.Where("id = ?", req.ID).First(&book).Error
	if err != nil {
		return nil, err
	}

	// 准备更新字段
	updates := map[string]interface{}{
		"title": req.Title,
		"count": req.Count,
		"isbn":  req.ISBN,

		"author":  req.Author,
		"content": req.Content,
		"summary": req.Summary,

		"version": book.Version + 1,
	}

	// 使用乐观锁更新
	result := d.db.Model(&model.Book{}).
		Where("id = ? AND version = ?", req.ID, book.Version).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("更新失败：数据已被其他用户修改，请刷新后重试")
	}

	// 重新查询更新后的数据
	err = d.db.Where("id = ?", req.ID).First(&book).Error
	if err != nil {
		return nil, err
	}

	return &book, nil
}

// BookListDAO 支持分页的书籍列表查询
// TODO 深分页问题
func (d *dbService) BookListDAO(req *api.BookSearchReq) (*api.BookSearchResp, error) {
	dbSql := d.db.Model(&model.Book{})
	if req.Title != "" {
		dbSql = dbSql.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.ISBN != "" {
		dbSql = dbSql.Where("isbn = ?", req.ISBN)
	}

	if req.Author != "" {
		dbSql = dbSql.Where("author LIKE ?", "%"+req.Author+"%")
	}
	if req.Content != "" {
		dbSql = dbSql.Where("content LIKE ?", "%"+req.Content+"%")
	}

	// 获取总数（用于分页）
	var total int64
	if err := dbSql.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count books: %w", err)
	}

	// 分页处理
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 4. 查询分页数据
	var books []model.Book
	if err := dbSql.Offset(offset).Limit(pageSize).Find(&books).Error; err != nil {
		return nil, fmt.Errorf("failed to query books: %w", err)
	}

	// 5. 转换为 API 响应结构
	bookResps := make([]api.BookInfoResp, 0, len(books))
	for _, book := range books {
		bookResps = append(bookResps, api.BookInfoResp{
			ID:      book.ID,
			Title:   book.Title,
			Author:  book.Author,
			Count:   book.Count,
			ISBN:    book.ISBN,
			Summary: book.Summary, // 假设 model.Book 有 Summary 字段
		})
	}

	// 6. 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &api.BookSearchResp{
		Books:      bookResps,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// BookGetByIDDAO 根据ID获取书籍详情
func (d *dbService) BookGetByIDDAO(id uint) (*model.Book, error) {
	var book model.Book
	err := d.db.Where("id = ?", id).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// BookGetByISBNDAO 根据ISBN获取书籍
func (d *dbService) BookGetByISBNDAO(isbn string) (*model.Book, error) {
	var book model.Book
	err := d.db.Where("isbn = ?", isbn).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}
