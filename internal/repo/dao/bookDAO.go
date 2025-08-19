package dao

import (
	"LibraryManagement/internal/api"
	"LibraryManagement/internal/model"
	"errors"
	"strconv"
)

type bookDAO interface {
	BookAddDAO(req *api.BookInfoReq) error
	BookDeleteDAO(idStr []string) error
	BookUpdateDAO(req *api.BookUpdateReq) error
	BookListDAO(req *api.BookSearchReq) (*[]api.BookInfoResp, error)
}

func (d *dbService) BookAddDAO(req *api.BookInfoReq) error {
	book := &model.Book{
		Title: req.Title,
		Count: req.Count,
		ISBN:  req.ISBN,
	}

	return d.db.Create(book).Error
}

func (d *dbService) BookDeleteDAO(idStr []string) error {
	ids := make([]uint, 0, len(idStr))
	for _, id := range idStr {
		uid, _ := strconv.ParseUint(id, 10, 64)
		ids = append(ids, uint(uid))
	}
	return d.db.Delete(&model.Book{}, ids).Error
}

// 数据库乐观锁更新
func (d *dbService) BookUpdateDAO(req *api.BookUpdateReq) error {
	var book model.Book
	err := d.db.Where("id = ?", req.ID).First(&book).Error
	if err != nil {
		return err
	}

	// 准备更新字段
	updates := map[string]interface{}{
		"title":   req.Title,
		"count":   req.Count,
		"isbn":    req.ISBN,
		"version": book.Version + 1,
	}

	// 使用乐观锁更新
	result := d.db.Model(&model.Book{}).
		Where("id = ? AND version = ?", req.ID, book.Version).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("更新失败：数据已被其他用户修改，请刷新后重试")
	}

	return nil
}

// TODO 考虑分页
func (d *dbService) BookListDAO(req *api.BookSearchReq) (*[]api.BookInfoResp, error) {
	dbSql := d.db.Model(&model.Book{})
	if req.Title != "" {
		dbSql = dbSql.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.ISBN != "" {
		dbSql = dbSql.Where("isbn = ?", req.ISBN)
	}
	var books []api.BookInfoResp
	dbSql.Find(&books)
	return &books, dbSql.Error
}
