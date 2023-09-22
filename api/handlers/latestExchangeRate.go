package handlers

import (
	"encoding/json"
	"fmt"
	"go-rich/models"
	"log"
	"net/http"
	"os"
	"time"
	"go-rich/pubsub/utils"

	"github.com/joho/godotenv"
)

func LatestExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	apiKey := os.Getenv("API_KEY")
	currency := r.URL.Query().Get("currency")
	apiEndpoint := fmt.Sprintf("https://api.currencyfreaks.com/latest?apikey=%s", apiKey)

	message := models.ExchangeRateMessage{
		Currency: currency,
		Url:      apiEndpoint,
	}

	utils.Publish(message, "amqp://localhost:5672", "exchange_rates")

	msgs, ch := utils.Consumer("amqp://localhost:5672", "api")
	defer ch.Close()

	var result models.ExchangeRateResult

	go func() {
		for d := range msgs {
			err := json.Unmarshal(d.Body, &result)
			if err != nil {
				panic(err)
			}
		}
	}()

	time.Sleep(3 * time.Second)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
