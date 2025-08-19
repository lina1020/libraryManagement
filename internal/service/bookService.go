package service

import (
	"LibraryManagement/internal/api"
	"LibraryManagement/internal/repo/dao"
)

type BookService interface {
	Add(dto *api.BookInfoReq) error
	Delete(ids []string) error
	Update(dto *api.BookUpdateReq) error
	List(dto *api.BookSearchReq) (*[]api.BookInfoResp, error)
}

type bookServiceImpl struct {
}

func NewBookService() BookService {
	return &bookServiceImpl{}
}

func (b *bookServiceImpl) Add(dto *api.BookInfoReq) error {

	err := dao.ApiDao.BookAddDAO(dto)

	return err

}

func (b *bookServiceImpl) Delete(ids []string) error {
	err := dao.ApiDao.BookDeleteDAO(ids)

	return err

}

// 数据库乐观锁更新
func (b *bookServiceImpl) Update(dto *api.BookUpdateReq) error {

	err := dao.ApiDao.BookUpdateDAO(dto)

	return err
}

func (b *bookServiceImpl) List(dto *api.BookSearchReq) (*[]api.BookInfoResp, error) {

	books, err := dao.ApiDao.BookListDAO(dto)

	return books, err
}
