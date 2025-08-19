package dao

import (
	"LibraryManagement/internal/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	ApiDao ApiDBDao
)

type dbService struct {
	db *gorm.DB
}

// GetDB returns the global db instance
func GetDB() *gorm.DB {
	return ApiDao.(*dbService).db
}

type ApiDBDao interface {
	bookDAO
	userDAO
}

func SetupDBLink() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config.Db.User,
		config.Config.Db.Password,
		config.Config.Db.Host,
		config.Config.Db.Port,
		config.Config.Db.Db,
	)
	fmt.Println("DSN:", dsn)

	var err error
	s := &dbService{}
	s.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取底层 *sql.DB
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")

	ApiDao = s

	return nil
}
