package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go-rich/models"
	mb "go-rich/pubsub/message_broker"
)

var result models.ExchangeRateResult

func main() {

	rabbitmq, err := mb.NewRabbitMQ("amqp://localhost:5672")
	if err != nil {
		panic(err)
	}

	// msgs, ch := utils.Consume("amqp://localhost:5672", "exchange_rates")

	msgs, err := rabbitmq.ReceiveMessage("exchange_rates")
	if err != nil {
		panic(err)
	}

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

			fmt.Println("RESULT",result)

			rabbitmq.SendMessage(result, "api")
			d.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
