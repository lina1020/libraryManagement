package controller

import (
	"LibraryManagement/config"
	"LibraryManagement/dto"
	"LibraryManagement/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

// AddBookController 添加书籍
func AddBookController(c *gin.Context) {
	bookInfoDTO := &dto.BookInfoDTO{}
	err := c.BindJSON(bookInfoDTO)
	if err != nil {
		return
	}
	fmt.Println("收到请求---bookAdd: ", bookInfoDTO)
	service.NewBookService(c, config.DB).Add(bookInfoDTO)
}

// DeleteBookController 删除书籍
func DeleteBookController(c *gin.Context) {
	var ids = make([]string, 0)
	ids = c.QueryArray("ids")

	fmt.Println("收到请求---Delete: ", ids)
	service.NewBookService(c, config.DB).Delete(ids)
}

// UpdateBookController 更新书籍
func UpdateBookController(c *gin.Context) {
	bookUpdateDTO := &dto.BookUpdateDTO{}
	err := c.BindJSON(bookUpdateDTO)
	if err != nil {
		return
	}
	fmt.Println("收到请求---bookUpdateDTO: ", bookUpdateDTO)
	service.NewBookService(c, config.DB).Update(bookUpdateDTO)
}

// BookListController 批量查询
func BookListController(c *gin.Context) {
	bookSearchDTO := &dto.BookSearchDTO{}
	err := c.BindJSON(bookSearchDTO)
	if err != nil {
		return
	}
	fmt.Println("收到请求---bookSearchDTO: ", bookSearchDTO)
	service.NewBookService(c, config.DB).List(bookSearchDTO)
}
