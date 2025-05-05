package config

import (
	"log"

	"gorm.io/driver/mysql" // Thay đổi driver từ postgres sang mysql
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Chuỗi kết nối MySQL (cập nhật thông tin user, password, và dbname)
	dsn := "root:my-secret-pw@tcp(127.0.0.1:3306)/task_management?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = database
}
