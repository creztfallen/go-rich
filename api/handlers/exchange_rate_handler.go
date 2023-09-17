package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ExchangeRateResponse struct {
	Date  string            `json:"date"`
	Base  string            `json:"base"`
	Rates map[string]string `json:"rates"`
}

type CurrencyRateResponse struct {
	Date     string `json:"date"`
	Base     string `json:"base"`
	Currency string `json:"rates"`
}

func LatestExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	apiKey := os.Getenv("API_KEY")
	currency := r.URL.Query().Get("currency")
	apiEndpoint := fmt.Sprintf("https://api.currencyfreaks.com/latest?apikey=%s", apiKey)

	response, err := http.Get(apiEndpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("API request failed with status code %d", response.StatusCode), http.StatusInternalServerError)
		return
	}

	var exchangeRateResponse ExchangeRateResponse
	err = json.NewDecoder(response.Body).Decode(&exchangeRateResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if currency != "" {
		base := exchangeRateResponse.Base
		date := exchangeRateResponse.Date
		rate, exists := exchangeRateResponse.Rates[currency]
		if !exists {
			http.Error(w, fmt.Sprintf("Currency %s not supported", currency), http.StatusNotFound)
			return
		}

		responseData := CurrencyRateResponse{
			Date:     date,
			Base:     base,
			Currency: rate,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exchangeRateResponse)
}

func HistoricalExchangeRateHandler(w http.ResponseWriter, r *http.Request) {

}
