package router

import (
	"LibraryManagement/internal/handler"
	"LibraryManagement/internal/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter(bookHandler *handler.BookHandler, userHandler *handler.UserHandler) *gin.Engine {
	router := gin.Default()

	register(router, bookHandler, userHandler)

	return router
}

func register(router *gin.Engine, bookHandler *handler.BookHandler, userHandler *handler.UserHandler) {

	// 公共路由（无需认证）
	auth := router.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	// 受保护路由
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware("")) // 所有登录用户可访问
	{
		api.POST("/books/list", bookHandler.BookList)

		api.GET("/books/:id", bookHandler.GetBook) // 获取单本书籍详情

		// ES搜索功能
		api.POST("/books/search", bookHandler.SearchBooks)            // 综合搜索
		api.GET("/books/search/title", bookHandler.SearchByTitle)     // 标题搜索
		api.GET("/books/search/content", bookHandler.SearchByContent) // 内容搜索
	}

	// 管理员专用路由
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware("admin")) // 仅管理员
	{
		admin.POST("/books/add", bookHandler.AddBook)
		admin.PUT("/books/update", bookHandler.UpdateBook)
		admin.DELETE("/books/delete", bookHandler.DeleteBook)

		// ES索引管理
		admin.POST("/es/index/init", bookHandler.InitESIndex)     // 初始化ES索引
		admin.POST("/es/index/reindex", bookHandler.ReindexBooks) // 重新索引所有数据

	}

}
