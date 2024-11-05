package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB(dsn string) error {
	var err error
	// 打开数据库连接
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	// 在这里可以做数据库迁移等操作
	// db.AutoMigrate(&models.Product{})  // 可以在这里进行自动迁移

	return nil
}

func GetDB() *gorm.DB {
	return db
}
