package api

import (
	"fmt"
	"github.com/creztfallen/go-rich/api/handlers"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/latest", handlers.LatestExchangeRateHandler)
	http.HandleFunc("/historical", handlers.HistoricalExchangeRateHandler)

	port := "8080"

	fmt.Printf("Starting server on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
