package main

import (
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/schema"
)

type Todo struct {
	Id			int
	Title		string
	Description string
	CompletedAt	sql.NullString
	DeletedAt	sql.NullString
}

func main() {
	// Connection
	dsn := "achyut:achyut@tcp(localhost:3306)/gotodo?charset=utf8&parseTime=True&loc=Local"

	// GORM Config
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Create
	db.Create(&Todo{
		Title: "Test Title",
		Description: "Test Description",
	})

	// Retrieve
	var todoItem Todo
	db.Where("title like ?", "%New Title%").First(&todoItem)

	// Update
	db.Model(&todoItem).Update("title", "New Title")

	// Delete
	db.Delete(&todoItem)
}