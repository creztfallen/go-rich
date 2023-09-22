package utils

import (
	

	amqp "github.com/rabbitmq/amqp091-go"
)

func BrokerConfig(connectionUrl, queueName string) (q amqp.Queue, ch *amqp.Channel) {
	conn := Connection(connectionUrl)
	ch = OpenChannel(conn)
	q = DeclareQueue(ch, queueName)

	return q, ch
}

func Connection(connectionUrl string) *amqp.Connection {
	conn, err := amqp.Dial(connectionUrl)
	if err != nil {
		panic(err)
	}

	return conn
}

func DeclareQueue(ch *amqp.Channel, queueName string) amqp.Queue {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		panic(err)
	}

	return q
}

func OpenChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}

