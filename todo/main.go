package main

import (
	"log"
	"net/http"
	"html/template"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type TodoItem struct {
	Id			int
	Title		string
	Description	string
	Completed	bool
}

type PageData struct {
	Title string
	TodoItems []TodoItem
}

var db *sql.DB
var err error
var tmpl *template.Template
var databaseName = "gotodo"
var todoItems []TodoItem

func main() {
	db, err = sql.Open("mysql", "achyut:achyut@tcp(127.0.0.1:3306)/" + databaseName)

	if err != nil {
		log.Fatal("Unable to connect to DB: ", err)
	}
	defer db.Close()

	createDatabase()

	mux := http.NewServeMux()

	tmpl = template.Must(template.ParseFiles("static/index.gohtml"))
	if err != nil {
		log.Fatal("Unable to parse template: ", err)
	}
	mux.HandleFunc("/", indexHandler)

	log.Println("Listening on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func createDatabase() {
	createDatabase := `
		CREATE DATABASE IF NOT EXISTS gotodo;
	`
	_, err := db.Exec(createDatabase)
	if err != nil {
		log.Fatal("Unable to create database: ", err)
	}

	createTable := `
		CREATE TABLE IF NOT EXISTS todo (
			id int(11) NOT NULL AUTO_INCREMENT,
			title varchar(255) NOT NULL,
			description varchar(255) NOT NULL,
			completed boolean NOT NULL,
			PRIMARY KEY (id)
		)
		ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Unable to create table: ", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	todoItems = []TodoItem{}
	selectData := "SELECT * FROM " + databaseName + ".todo;"
	rows, err := db.Query(selectData)
	if err != nil {
		log.Fatal("Unable to select data: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		todoItem := TodoItem{}
		err := rows.Scan(&todoItem.Id, &todoItem.Title, &todoItem.Description, &todoItem.Completed)
		if err != nil {
			log.Fatal("Unable to scan data: ", err)
		}
		todoItems = append(todoItems, todoItem)
	}

	data := PageData {
		Title: "Todo",
		TodoItems: todoItems,
	}
	tmpl.Execute(w, data)
}