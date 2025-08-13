package dao

import (
	"LibraryManagement/dto"
	"LibraryManagement/entity"
	"LibraryManagement/vo"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

func BookAddDAO(db *gorm.DB, dto *dto.BookInfoDTO) error {
	book := &entity.Book{
		Title: dto.Title,
		Count: dto.Count,
		ISBN:  dto.ISBN,
	}

	return db.Create(book).Error
}

func BookDeleteDAO(db *gorm.DB, idStr []string) error {
	ids := make([]uint, 0, len(idStr))
	for _, id := range idStr {
		uid, _ := strconv.ParseUint(id, 10, 64)
		ids = append(ids, uint(uid))
	}
	return db.Delete(&entity.Book{}, ids).Error
}

// 数据库乐观锁更新
func BookUpdateDAO(db *gorm.DB, dto *dto.BookUpdateDTO) error {
	var book entity.Book
	err := db.Where("id = ?", dto.ID).First(&book).Error
	if err != nil {
		return err
	}

	// 准备更新字段
	updates := map[string]interface{}{
		"title":   dto.Title,
		"count":   dto.Count,
		"isbn":    dto.ISBN,
		"version": book.Version + 1,
	}

	// 使用乐观锁更新
	result := db.Model(&entity.Book{}).
		Where("id = ? AND version = ?", dto.ID, book.Version).
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
func BookListDAO(db *gorm.DB, dto *dto.BookSearchDTO) (*[]vo.BookInfoVO, error) {
	dbSql := db.Model(&entity.Book{})
	if dto.Title != "" {
		dbSql = dbSql.Where("title LIKE ?", "%"+dto.Title+"%")
	}
	if dto.ISBN != "" {
		dbSql = dbSql.Where("isbn = ?", dto.ISBN)
	}
	var books []vo.BookInfoVO
	dbSql.Find(&books)
	return &books, dbSql.Error
}
