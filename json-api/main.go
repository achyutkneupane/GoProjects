package main

import (
	"log"
	"encoding/json"
	"net/http"
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/schema"
	"time"
)

type Todo struct {
	Id			int
	Title		string
	Description string
	CompletedAt	sql.NullString
	DeletedAt	sql.NullString
}

var databaseName = "gotodo"
var db *gorm.DB
var err error
var mux *http.ServeMux

func main() {

	dsn := "achyut:achyut@tcp(localhost:3306)/gotodo?charset=utf8&parseTime=True&loc=Local"

	// GORM Config
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect database")
	}

	mux = http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/add", addTodoHandler)
	mux.HandleFunc("/completed", completedHandler)
	mux.HandleFunc("/deleted", deletedHandler)

	log.Println("Listening on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var todoItems []Todo
	db.Find(&todoItems)
	json.NewEncoder(w).Encode(todoItems)
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	if(r.Method == "POST") {
		var todo Todo
		todo.Title = r.FormValue("title")
		todo.Description = r.FormValue("description")
		db.Create(&todo)
		json.NewEncoder(w).Encode(todo)
	} else {
		json.NewEncoder(w).Encode("Method not allowed")
	}
}

func completedHandler(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	id := r.FormValue("id")
	db.First(&todo, id)
	todo.CompletedAt = sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true}
	db.Save(&todo)
	json.NewEncoder(w).Encode(todo)
}

func deletedHandler(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	id := r.FormValue("id")
	db.First(&todo, id)
	todo.DeletedAt = sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true}
	db.Save(&todo)
	json.NewEncoder(w).Encode(todo)
}