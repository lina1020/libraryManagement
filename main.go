package main

import (
	"LibraryManagement/config"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// 初始化配置，连接数据库
func init() {
	if err := config.LoadConfig("./config.yaml"); err != nil {
		log.Fatal(err)
	}

	err := config.SetupDBLink()
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	engine := gin.Default()
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	err := engine.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
