package main

import (
	"LibraryManagement/internal/config"
	"LibraryManagement/internal/es"
	"LibraryManagement/internal/handler"
	"LibraryManagement/internal/repo/dao"
	"LibraryManagement/internal/router"
	"LibraryManagement/internal/service"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Gin 框架 Web 服务的优雅关闭（Graceful Shutdown）实现
func main() {

	// 初始化配置，连接数据库
	if err := config.LoadConfig("./config.yaml"); err != nil {
		log.Fatal("加载配置文件失败: ", err)
	}

	// 初始化数据库连接
	if err := dao.SetupDBLink(); err != nil {
		log.Fatal("数据库连接失败: ", err)
	}

	// ES初始化
	if err := es.InitES(); err != nil {
		log.Printf("ES 初始化失败: %v", err)
		log.Println("将在没有ES支持的情况下继续运行")
	} else {
		log.Println("ES 初始化成功")
	}

	// 依赖注入
	// init service
	bookService := service.NewBookService()
	userService := service.NewUserService()

	// 初始化ES索引（如果ES可用）
	if es.Client != nil {
		if err := bookService.InitializeESIndex(); err != nil {
			log.Printf("ES索引初始化失败: %v", err)
		} else {
			log.Println("ES索引初始化成功")
		}
	}

	// init handler
	bookHandler := handler.NewBookHandler(bookService)
	userHandler := handler.NewUserHandler(userService)

	gin := router.InitRouter(bookHandler, userHandler)

	//创建HTTP服务器
	server := &http.Server{
		Addr:    config.Config.Server.Port,
		Handler: gin,
	}

	//启动HTTP服务器
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	//等待退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	//创建超时上下文，Shutdown可以让未处理的连接在这个时间内关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//停止HTTP服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
