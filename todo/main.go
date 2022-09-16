package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type TodoItem struct {
	Id			int
	Title		string
	Description string
	Completed	bool
	DeletedAt	sql.NullString
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
var taskId string

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
	mux.HandleFunc("/add", addTodoHandler)
	mux.HandleFunc("/completed", completedHandler)

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
			description varchar(255) NULL,
			completed boolean NOT NULL,
			deleted_at datetime NULL,
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
	selectData := "SELECT * FROM " + databaseName + ".todo ORDER BY completed ASC, title ASC;"
	rows, err := db.Query(selectData)
	if err != nil {
		log.Fatal("Unable to select data: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		todoItem := TodoItem{}
		err := rows.Scan(&todoItem.Id, &todoItem.Title, &todoItem.Description, &todoItem.Completed, &todoItem.DeletedAt)
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

func completedHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/completed" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		taskId = r.FormValue("id")
		completeTask := "UPDATE todo SET completed=1 WHERE id=" + taskId + ";"
		_, err = db.Exec(completeTask)
		if err != nil {
			log.Fatal("Unable to complete task: ", err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/add" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		title := r.FormValue("title")
		description := r.FormValue("description")
		insertData := "INSERT INTO todo (title, description, completed) VALUES (?, ?, ?);"
		_, err = db.Exec(insertData, title, description, 0)
		if err != nil {
			log.Fatal("Unable to insert data: ", err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		if(r.Method == "GET") {
			http.ServeFile(w, r, "static/insert.html")
		} else {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		}
	}
}