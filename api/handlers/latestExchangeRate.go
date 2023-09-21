package handlers

import (
	"encoding/json"
	"fmt"
	"go-rich/models"
	"go-rich/publishers"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
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

	publishers.Publish(message)

	conn, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	q, err := ch.QueueDeclare(
		"api", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	var result models.ExchangeRateResult

	go func() {
		for d := range msgs {
			err = json.Unmarshal(d.Body, &result)
		}
	}()

	time.Sleep(3 * time.Second)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func HistoricalExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
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

	publishers.Publish(message)

	conn, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	q, err := ch.QueueDeclare(
		"api", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	var result models.ExchangeRateResult

	go func() {
		for d := range msgs {
			err = json.Unmarshal(d.Body, &result)
		}
	}()

	time.Sleep(3 * time.Second)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}