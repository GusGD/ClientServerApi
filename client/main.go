package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	APIURLLOCAL = "http://localhost:8080/"
	TimeoutReq  = 30000 * time.Millisecond
)

var client = &http.Client{}

type Currency struct {
	Bid string `json:"bid"`
}

func main() {
	err := getCurrency()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}

func getCurrency() error {
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutReq)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", APIURLLOCAL, nil)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Error: Execution time exceeded")
		}
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var data map[string]Currency
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	bidString := data["USDBRL"].Bid

	err = writeToFile("cotacao.txt", bidString)
	if err != nil {
		return err
	}
	fmt.Println("A cotação atual do dólar é R$:", bidString)
	return nil
}

func writeToFile(fileName, content string) error {
	err := os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Valor:%s", content)), 0644)
	if err != nil {
		return err
	}
	return nil
}
