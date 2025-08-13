package service

import (
	"LibraryManagement/dao"
	"LibraryManagement/dto"
	"LibraryManagement/result"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type BookService interface {
	Add(dto *dto.BookInfoDTO)
	Delete(ids []string)
	Update(dto *dto.BookUpdateDTO)
	List(dto *dto.BookSearchDTO)
}

type bookServiceImpl struct {
	Ctx *gin.Context
	Db  *gorm.DB
}

func (b *bookServiceImpl) Add(dto *dto.BookInfoDTO) {
	// 验证器会根据结构体里写的 validate 标签，自动检查字段是否符合规则
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(b.Ctx, result.RequiredCode, result.GetMessage(result.RequiredCode))
		return
	}
	// TODO ISBN校验

	err = dao.BookAddDAO(b.Db, dto)
	if err != nil {
		result.Failed(b.Ctx, result.FailedCode, "书籍添加失败:"+err.Error())
		return
	}

	result.Success(b.Ctx, "书籍添加成功")

}

func (b *bookServiceImpl) Delete(ids []string) {
	err := dao.BookDeleteDAO(b.Db, ids)
	if err != nil {
		result.Failed(b.Ctx, result.FailedCode, "书籍删除失败:"+err.Error())
		return
	}

	result.Success(b.Ctx, "书籍删除成功")

}

// 数据库乐观锁更新
func (b *bookServiceImpl) Update(dto *dto.BookUpdateDTO) {
	// 验证器会根据结构体里写的 validate 标签，自动检查字段是否符合规则
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(b.Ctx, result.RequiredCode, result.GetMessage(result.RequiredCode))
		return
	}

	err = dao.BookUpdateDAO(b.Db, dto)
	if err != nil {
		result.Failed(b.Ctx, result.FailedCode, "书籍更新失败:"+err.Error())
		return
	}

	result.Success(b.Ctx, "书籍更新成功")
}

func (b *bookServiceImpl) List(dto *dto.BookSearchDTO) {
	// 验证器会根据结构体里写的 validate 标签，自动检查字段是否符合规则
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(b.Ctx, result.RequiredCode, result.GetMessage(result.RequiredCode))
		return
	}

	books, err := dao.BookListDAO(b.Db, dto)
	if err != nil {
		result.Failed(b.Ctx, result.FailedCode, "书籍查询失败:"+err.Error())
		return
	}

	result.Success(b.Ctx, books)
}

func NewBookService(ctx *gin.Context, db *gorm.DB) BookService {
	return &bookServiceImpl{
		Ctx: ctx,
		Db:  db,
	}
}
