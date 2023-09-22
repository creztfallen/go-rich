package utils

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(message interface{}, connectionString, queueName string) {
	q, ch := BrokerConfig(connectionString, queueName)
	amqpMessage, ctx, cancel := CreateMessage(message)
	defer cancel()

	PublishMessage(ch, amqpMessage, q, ctx)
}

func CreateMessage(message interface{}) (amqp.Publishing, context.Context, context.CancelFunc){
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	amqpMessage := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(jsonMessage),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	return amqpMessage, ctx, cancel
}

func PublishMessage(ch *amqp.Channel, amqpMessage amqp.Publishing, q amqp.Queue, ctx context.Context) {
	err := ch.PublishWithContext(ctx,
		"",          // exchange
		q.Name,      // routing key
		false,       // mandatory
		false,       // immediate
		amqpMessage, // message
	)
	if err != nil {
		panic(err)
	}
}