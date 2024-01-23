package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Currency struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/", nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	out, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data map[string]Currency
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	bidString := data["USDBRL"].Bid

	err = os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Valor:%s", bidString)), 0644)
	if err != nil {
		panic(err)
	}

	defer out.Close()
	fmt.Println("A cotação atual do dólar é R$:", bidString)
}
