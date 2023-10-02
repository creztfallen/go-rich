package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go-rich/models"
	mb "go-rich/pubsub/message_broker"

	"github.com/joho/godotenv"
)

var result models.ExchangeRateResult
var exchangeRateResponse models.ExchangeRateResponse
var message models.ExchangeRateMessage

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	connection_uri := os.Getenv("CONNECTION_URI")
	messageBroker, err := mb.NewRabbitMQ(connection_uri)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msgs, err := messageBroker.ReceiveMessage("exchange_rates")
	if err != nil {
		panic(err)
	}

	delivery := mb.ConvertToAMQPDeliveryChan(msgs)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-delivery:
				if !ok {
					return
				}

				fmt.Println(string(d.Body))
				err := json.Unmarshal(d.Body, &message)
				if err != nil {
					log.Printf("Error unmarshalling message: %s", err)
					continue
				}

				go func() {
					response, err := http.Get(message.Url)
					if err != nil {
						log.Printf("Error getting exchange rate: %s", err)
						return
					}

					err = json.NewDecoder(response.Body).Decode(&exchangeRateResponse)
					if err != nil {
						log.Printf("Error decoding exchange rate response: %s", err)
						return
					}

					result = models.ExchangeRateResult{
						Date: exchangeRateResponse.Date,
						Base: exchangeRateResponse.Base,
						Rate: exchangeRateResponse.Rates[message.Currency],
					}

					fmt.Println("RESULT", result)

					err = messageBroker.SendMessage(result, "api")
					if err != nil {
						log.Printf("Error sending message: %s", err)
					}
				}()
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	// Wait for SIGINT or SIGTERM signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Cancel context to stop message processing
	cancel()
}
