package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Currency struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/", HomeHandle)
	http.ListenAndServe(":8080", nil)
}

func HomeHandle(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "../currencies.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS currencies(id INTEGER PRIMARY KEY, bid TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Debug:", string(body))
	var data map[string]Currency
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	bidString := data["USDBRL"].Bid
	fmt.Println("bidString")
	_, err = db.Exec("INSERT INTO currencies (bid) VALUES (?)", bidString)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(data)
}
