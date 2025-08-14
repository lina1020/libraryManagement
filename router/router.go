package router

import (
	"LibraryManagement/controller"
	"LibraryManagement/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	router := gin.Default()

	register(router)

	return router
}

func register(router *gin.Engine) {

	// 公共路由（无需认证）
	auth := router.Group("/auth")
	{
		auth.POST("/register", controller.RegisterController)
		auth.POST("/login", controller.LoginController)
	}

	// 受保护路由
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware("")) // 所有登录用户可访问
	{
		api.GET("/books/list", controller.BookListController)
		// TODO 借阅归还等
	}

	// 管理员专用路由
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware("admin")) // 仅管理员
	{
		admin.POST("/books/add", controller.AddBookController)
		admin.PUT("/books/update", controller.DeleteBookController)
		admin.DELETE("/books/delete", controller.UpdateBookController)
	}

	//book := router.Group("/books")
	//{
	//	book.POST("/add", controller.AddBookController)
	//	book.DELETE("/delete", controller.DeleteBookController)
	//	book.PUT("/update", controller.UpdateBookController)
	//	book.GET("/list", controller.BookListController)
	//
	//}
}
