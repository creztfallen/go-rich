package main

import (
	"fmt"
	"net/http"
	"go-rich/api/handlers"
	mb "go-rich/pubsub/message_broker"
)


func main() {
	rabbitmq, err := mb.NewRabbitMQ("amqp://localhost:5672")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/latest", handlers.LatestExchangeRateHandler(*rabbitmq))
	http.HandleFunc("/latests", handlers.LatestExchangeRatesHandler)

	port := "8080"

	fmt.Printf("Starting server on port %s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
