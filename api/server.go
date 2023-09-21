package main

import (
	"fmt"
	"net/http"
	"go-rich/api/handlers"
)

func main() {
	http.HandleFunc("/latest", handlers.LatestExchangeRateHandler)

	port := "8080"

	fmt.Printf("Starting server on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
