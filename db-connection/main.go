package main

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var id int
	var name string
	var email string

	databaseName := "gotest"
	db, err := sql.Open("mysql", "achyut:achyut@tcp(127.0.0.1:3306)/" + databaseName)
	if err != nil {
		log.Fatal("Unable to connect to DB: ", err)
	}
	fmt.Println("Connected to database " + databaseName)
	defer db.Close()

	// SQL command for creating `user` table
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id int(11) NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		email varchar(255) NOT NULL,
		PRIMARY KEY (id)
	)
	ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`
	result, err := db.Exec(createTable)
	if err != nil {
		log.Fatal("Unable to create table: ", err)
	}
	fmt.Println("Table created successfully: ", result)

	// SQL command for inserting data into `user` table
	insertData := `
		INSERT INTO users (name, email) VALUE
		('Achyut', 'achyutkneupane@gmail.com');
	`
	result, err = db.Exec(insertData)
	if err != nil {
		log.Fatal("Unable to insert data: ", err)
	}
	fmt.Println("Data inserted successfully: ", result)

	// SQL command for selecting data from `user` table
	selectData := `
		SELECT * FROM users;
	`
	rows, err := db.Query(selectData)
	if err != nil {
		log.Fatal("Unable to select data: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name, &email)
		if err != nil {
			log.Fatal("Unable to scan data: ", err)
		}
		fmt.Println("\nID:\t", id, "\nname:\t", name, "\nemail:\t", email)
	}
}