package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./currencies.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var coluna1 int
	var coluna2, coluna3 string

	err = db.QueryRow("SELECT * FROM currencies").Scan(&coluna1, &coluna2, &coluna3)
	if err != nil {
		panic(err)
	}

	fmt.Println(coluna1, coluna2)
}
