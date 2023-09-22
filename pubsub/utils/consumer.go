package utils

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func Consumer(connectionString, queueName string) (<-chan amqp.Delivery, *amqp.Channel) {

	q, ch := BrokerConfig(connectionString, queueName)

	DeclareQos(ch, 1, 0, false)

	msgs:= Consume(ch, q)

	return msgs, ch
}

func DeclareQos(ch *amqp.Channel, prefetchCount, prefetchSize int, global bool) {
	err := ch.Qos(
		prefetchCount, // prefetch count
		prefetchSize,  // prefetch size
		global,        // global
	)
	if err != nil {
		panic(err)
	}
}

func Consume(ch *amqp.Channel, q amqp.Queue) <- chan amqp.Delivery{
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

	return msgs
}