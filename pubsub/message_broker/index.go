package message_broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
    SendMessage(message interface{}, queueName string) error
    ReceiveMessage(queueName string) (<- chan amqp.Delivery, error)
}

type RabbitMQ struct {
	conn *amqp.Connection
	ch *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	rabbitmqInstance := &RabbitMQ{
		conn: conn,
		ch: ch,
	}

	return rabbitmqInstance, nil
}