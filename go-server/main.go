package main

import (
	"fmt"
	"log"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		username := r.FormValue("username")
		fmt.Fprintf(w, "You are logged in as %s", username)
		
	} else {
		if r.Method == "GET" {
			http.ServeFile(w, r, "static/login.html")
		} else {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func main() {
	staticDirectory := http.Dir("./static")
	staticHandler := http.FileServer(staticDirectory)
	http.Handle("/", staticHandler)
	http.HandleFunc("/login", loginHandler)

	fmt.Println("Server running on port 8080")
	fmt.Println("Open http://localhost:8080 in your browser")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}