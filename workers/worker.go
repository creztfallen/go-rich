package main

import (
	"encoding/json"
	"fmt"
	"go-rich/models"
	"go-rich/publishers"
	"log"
	"net/http"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var result models.ExchangeRateResult

func main() {
	conn, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"exchange-rates", // name
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

			publishers.PublishBack(result)

		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
