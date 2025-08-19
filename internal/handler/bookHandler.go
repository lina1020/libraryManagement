package handler

import (
	"LibraryManagement/internal/api"
	"LibraryManagement/internal/api/result"
	"LibraryManagement/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BookHandler struct {
	bookService service.BookService
}

func NewBookHandler(bookService service.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

// AddBook 添加书籍
func (b *BookHandler) AddBook(c *gin.Context) {
	bookInfoReq := &api.BookInfoReq{}
	err := c.BindJSON(bookInfoReq)
	if err != nil {
		result.Failed(c, result.RequiredCode, result.GetMessage(result.RequiredCode))
		return
	}

	fmt.Println("收到请求---bookAdd: ", bookInfoReq)

	// 验证器会根据结构体里写的 validate 标签，自动检查字段是否符合规则
	err = validator.New().Struct(bookInfoReq)
	if err != nil {
		result.Failed(c, result.RequiredCode, result.GetMessage(result.RequiredCode))
		return
	}
	// TODO ISBN校验

	err = b.bookService.Add(bookInfoReq)

	if err != nil {
		result.Failed(c, result.FailedCode, "书籍添加失败:"+err.Error())
		return
	}

	result.Success(c, "书籍添加成功")

}

// DeleteBook 删除书籍
func (b *BookHandler) DeleteBook(c *gin.Context) {
	var ids = make([]string, 0)
	ids = c.QueryArray("ids")
	fmt.Println("收到请求---Delete: ", ids)

	if len(ids) == 0 {
		result.Failed(c, result.RequiredCode, result.GetMessage(result.RequiredCode))
		return
	}

	err := b.bookService.Delete(ids)

	if err != nil {
		result.Failed(c, result.FailedCode, "书籍删除失败:"+err.Error())
		return
	}

	result.Success(c, "书籍删除成功")
}

// UpdateBook 更新书籍
func (b *BookHandler) UpdateBook(c *gin.Context) {
	bookUpdateReq := &api.BookUpdateReq{}
	err := c.BindJSON(bookUpdateReq)
	if err != nil {
		return
	}
	fmt.Println("收到请求---bookUpdateReq: ", bookUpdateReq)

	// 验证器会根据结构体里写的 validate 标签，自动检查字段是否符合规则
	err = validator.New().Struct(bookUpdateReq)
	if err != nil {
		result.Failed(c, result.RequiredCode, result.GetMessage(result.RequiredCode))
		return
	}

	err = b.bookService.Update(bookUpdateReq)
	if err != nil {
		result.Failed(c, result.FailedCode, "书籍更新失败:"+err.Error())
		return
	}
	result.Success(c, "书籍更新成功")
}

// BookList 批量查询
func (b *BookHandler) BookList(c *gin.Context) {
	bookSearchReq := &api.BookSearchReq{}
	err := c.BindJSON(bookSearchReq)
	if err != nil {
		return
	}
	fmt.Println("收到请求---bookSearchReq: ", bookSearchReq)

	// 验证器会根据结构体里写的 validate 标签，自动检查字段是否符合规则
	err = validator.New().Struct(bookSearchReq)
	if err != nil {
		result.Failed(c, result.RequiredCode, result.GetMessage(result.RequiredCode))
		return
	}
	books, err := b.bookService.List(bookSearchReq)
	if err != nil {
		result.Failed(c, result.FailedCode, "书籍查询失败:"+err.Error())
		return
	}

	result.Success(c, books)
}
