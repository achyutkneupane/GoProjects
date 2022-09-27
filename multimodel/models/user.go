package models

type User struct {
	Id			int `gorm:"primaryKey,autoIncrement,not null"`
	Username	string `gorm:"unique,not null"`
	Password	string `gorm:"not null"`
	Email		string `gorm:"unique,not null"`
	FirstName	string `gorm:"not null"`
	LastName	string `gorm:"not null"`
}