package router

import (
	"LibraryManagement/controller"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	router := gin.Default()

	register(router)

	return router
}

func register(router *gin.Engine) {

	book := router.Group("/books")
	{
		book.POST("/add", controller.AddBookController)
		book.DELETE("/delete", controller.DeleteBookController)
		book.PUT("/update", controller.UpdateBookController)
		book.GET("/list", controller.BookListController)

	}
}
