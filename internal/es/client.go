package es

import (
	"LibraryManagement/internal/config"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var Client *elasticsearch.Client

func InitES() error {
	if config.Config.Elasticsearch.Host == "" {
		log.Println("Elasticsearch配置为空，跳过ES客户端初始化")
		return nil
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%d", config.Config.Elasticsearch.Host, config.Config.Elasticsearch.Port),
		},
	}

	if config.Config.Elasticsearch.Username != "" {
		cfg.Username = config.Config.Elasticsearch.Username
		cfg.Password = config.Config.Elasticsearch.Password
	}

	var err error
	Client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create ES client: %w", err)
	}

	// 测试连接
	res, err := Client.Info()
	if err != nil {
		return fmt.Errorf("failed to get ES info: %w", err)
	}
	defer res.Body.Close()

	log.Println("Successfully connected to Elasticsearch")
	return nil
}
