package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error
var dbHost, dbPort, dbUser, dbPass, dbName string
var dsn string

type User struct {
	Id        int    `gorm:"primaryKey,autoIncrement,not null"`
	Username  string `gorm:"unique,not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"unique,not null"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Gender    string `gorm:"not null"`
	CreatedAt string
	DeletedAt sql.NullString
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbUser = os.Getenv("DB_USERNAME")
	dbPass = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")

	dsn = dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8&parseTime=True&loc=Local"

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})

	mux := http.NewServeMux()
	handleRoutes(mux)
	fmt.Println("Server is running on url http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func handleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/logout", logoutHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	// check if cookie exist
	cookie, _ := r.Cookie("goauthsession")

	if cookie != nil {
		// get user from database
		var data User
		db.Where("username = ?", cookie.Value).First(&data)

		if data.Username == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if data.Username != "" {
			tmpl := template.Must(template.ParseFiles("static/dashboard.html"))
			tmpl.Execute(w, data)
		}
	} else {
		tmpl := template.Must(template.ParseFiles("static/index.html"))
		tmpl.Execute(w, nil)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	cookie, _ := r.Cookie("goauthsession")
	if cookie != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		password = fmt.Sprintf("%x", sha1.Sum([]byte(password)))

		var user User
		db.Where("username = ? AND password = ?", username, password).First(&user)

		if user.Username == "" {
			http.Error(w, "Invalid username or password", http.StatusForbidden)
			return
		}

		// Set cookie
		cookie := http.Cookie{
			Name:    "goauthsession",
			Value:   username,
			Expires: time.Now().Add(24 * time.Hour),
		}

		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		tmpl := template.Must(template.ParseFiles("static/login.html"))
		tmpl.Execute(w, nil)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	cookie, _ := r.Cookie("goauthsession")
	if cookie != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		gender := r.FormValue("gender")
		password := r.FormValue("password")
		password = fmt.Sprintf("%x", sha1.Sum([]byte(password)))

		user := User{
			Username:  username,
			Password:  password,
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
			Gender:    gender,
			CreatedAt: time.Now().String()[:19],
		}

		db.Create(&user)

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	} else {
		tmpl := template.Must(template.ParseFiles("static/register.html"))
		tmpl.Execute(w, nil)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	cookie, err := r.Cookie("goauthsession")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
