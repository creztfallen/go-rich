package handlers

import (
	"encoding/json"
	"fmt"
	"go-rich/models"
	mb "go-rich/pubsub/message_broker"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LatestExchangeRateHandler(rabbitmq mb.MessageQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file %v", err)
		}

		apiUrl := os.Getenv("API_URL")
		currency := r.URL.Query().Get("currency")

		message := models.ExchangeRateMessage{
			Currency: currency,
			Url:      apiUrl,
		}

		rabbitmq.SendMessage(message, "exchange_rates")

		msgs, err := rabbitmq.ReceiveMessage("api")
		if err != nil {
			panic(err)
		}

		var result models.ExchangeRateResult
		var resultCh = make(chan models.ExchangeRateResult)

		go func() {
			for d := range msgs {
				err := json.Unmarshal(d.Body, &result)
				if err != nil {
					panic(err)
				}
				d.Ack(false)
				resultCh <- result
				fmt.Println("RESULT1", result)
			}
		}()
	

		select {
		case <-time.After(5 * time.Second):
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("Timeout waiting for result."))
		case result := <-resultCh:
			fmt.Println("RESULT2", result)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}
	}
}

func LatestExchangeRatesHandler(w http.ResponseWriter, r *http.Request) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	apiUrl := os.Getenv("API_URL")

	response, err := http.Get(apiUrl)
	if err != nil {
		panic(err)
	}

	var result models.ExchangeRateResponse

	json.NewDecoder(response.Body).Decode(&result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
