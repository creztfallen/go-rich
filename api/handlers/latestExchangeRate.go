package handlers

import (
	"encoding/json"
	"fmt"
	"go-rich/models"
	mb "go-rich/pubsub/message_broker"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func LatestExchangeRateHandler(w http.ResponseWriter, r *http.Request) {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	connection_uri := os.Getenv("CONNECTION_URI")
	apiUrl := os.Getenv("API_URL")
	currency := r.URL.Query().Get("currency")

	message := models.ExchangeRateMessage{
		Currency: currency,
		Url:      apiUrl,
	}

	rabbitmq, err := mb.NewRabbitMQ(connection_uri)
	if err != nil {
		panic(err)
	}
	defer rabbitmq.Close()

	rabbitmq.SendMessage(message, "exchange_rates")

	msgs, err := rabbitmq.ReceiveMessage("api")
	if err != nil {
		panic(err)
	}

	defer rabbitmq.Close()

	var result models.ExchangeRateResult
	var resultCh = make(chan models.ExchangeRateResult)

	go func() {
		for d := range msgs {
			err := json.Unmarshal(d.Body, &result)
			if err != nil {
				panic(err)
			}
			resultCh <- result
			fmt.Println("RESULT1", result)
		}
		close(resultCh)
	}()

	result = <-resultCh
	fmt.Println("RESULT2", result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
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
