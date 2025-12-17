package database

import (
	"firstRestAPI/todo"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=todo_pass dbname=todo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println("failed to connect database", err)
		return
	}

	err = db.AutoMigrate(&todo.Task{})
	if err != nil {
		fmt.Println(err)
		return
	}

	DB = db
}
