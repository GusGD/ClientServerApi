package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Currency struct {
	Coluna1 int
	Coluna2 string
	Coluna3 string
}

func main() {
	db, err := openDataBase("./currencies.db")
	if err != nil {
		log.Fatalf("Error Opening database: %v", err)
	}
	defer db.Close()

	currency, err := selectCurrency(db)
	if err != nil {
		log.Fatalf("Error Selecting Currency: %v", err)
	}
	fmt.Println(currency.Coluna1, currency.Coluna2, currency.Coluna3)
}

func openDataBase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func selectCurrency(db *sql.DB) (*Currency, error) {
	var currency Currency
	err := db.QueryRow("SELECT * FROM currencies").Scan(&currency.Coluna1, &currency.Coluna2, &currency.Coluna3)
	if err != nil {
		return nil, err
	}
	return &currency, nil
}
