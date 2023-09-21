package publishers

import (
	"context"
	"encoding/json"
	"go-rich/models"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)



func Publish(message models.ExchangeRateMessage) {
	conn, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
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
		panic(err)
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	amqpMessage := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(jsonMessage),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()


	err = ch.PublishWithContext(ctx,
		"",                 // exchange
		q.Name, 			// routing key
		false,              // mandatory
		false,              // immediate
		amqpMessage,        // message
	)
	if err != nil {
		panic(err)
	}
}

func PublishBack(message models.ExchangeRateResult) {
	conn, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"api", 			  // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		panic(err)
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	amqpMessage := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(jsonMessage),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()


	err = ch.PublishWithContext(ctx,
		"",                 // exchange
		q.Name, 			// routing key
		false,              // mandatory
		false,              // immediate
		amqpMessage,        // message
	)
	if err != nil {
		panic(err)
	}
}