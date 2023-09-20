package main

import (
	"fmt"
	"net/http"

	"go-rich/api/handlers"
)

var messageChannel = make(chan string)

func main() {
	http.HandleFunc("/latest", handlers.LatestExchangeRateHandler)
	http.HandleFunc("/historical", handlers.HistoricalExchangeRateHandler)

	port := "8080"

	fmt.Printf("Starting server on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
