package main

import (
	"encoding/json"

	"go-rich/models"
	"go-rich/pubsub/utils"

	"log"
	"net/http"
)

var result models.ExchangeRateResult

func main() {

	msgs, ch := utils.Consumer("amqp://localhost:5672", "exchange_rates")

	defer ch.Close()
	var forever chan struct{}

	go func() {
		for d := range msgs {
			var message models.ExchangeRateMessage
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				panic(err)
			}

			response, err := http.Get(message.Url)
			if err != nil {
				panic(err)
			}

			var exchangeRateResponse models.ExchangeRateResponse
			err = json.NewDecoder(response.Body).Decode(&exchangeRateResponse)
			if err != nil {
				panic(err)
			}

			result = models.ExchangeRateResult{
				Date: exchangeRateResponse.Date,
				Base: exchangeRateResponse.Base,
				Rate: exchangeRateResponse.Rates[message.Currency],
			}

			utils.Publish(result, "amqp://localhost:5672", "api")

		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
