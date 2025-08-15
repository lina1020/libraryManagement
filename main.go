package main

import (
	"LibraryManagement/config"
	"LibraryManagement/router"
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
		log.Fatal(err)
	}

	err := config.SetupDBLink()
	if err != nil {
		log.Fatal(err)
	}

	gin := router.InitRouter()

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
