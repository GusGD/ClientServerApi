package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DBPath        = "../currencies.db"
	APIURL        = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	TimeoutAPI    = 20000 * time.Millisecond
	TimeoutInsert = 100 * time.Millisecond
)

var client = &http.Client{}

type Currency struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/", HomeHandle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HomeHandle(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = createTable(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := fetchData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = insertData(db, data["USDBRL"].Bid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func createTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS currencies(id INTEGER PRIMARY KEY, bid TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)")
	return err
}

func fetchData() (map[string]Currency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutAPI)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", APIURL, nil)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Error: Execution time exceeded")
		}
		panic(err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var data map[string]Currency
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	return data, nil
}

func insertData(db *sql.DB, bid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutInsert)
	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO currencies (bid) VALUES (?)", bid)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Error: Execution time exceeded")
		}
	}
	return err
}
