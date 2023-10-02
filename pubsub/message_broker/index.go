package message_broker

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
    SendMessage(message interface{}, queueName string) error
    ReceiveMessage(queueName string) (<- chan interface{}, error)
	Close()
}

type RabbitMQ struct {
	conn *amqp.Connection
	ch *amqp.Channel
	CleanUp func()
}

type SQS struct {
	svc *sqs.SQS
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

	rabbitmqInstance.CleanUp = func() {
		ch.Close()
	}

	return rabbitmqInstance, nil
}

func NewSQS(svc *sqs.SQS) (*SQS, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	sqsClient := sqs.New(sess)

	awsQueueInstance := &SQS{
		svc: sqsClient,
	}

	return awsQueueInstance, nil
}

func ConvertToInterfaceChan(inputChan <-chan amqp.Delivery) <-chan interface{} {
    outputChan := make(chan interface{})
    
    go func() {
        defer close(outputChan)
        for msg := range inputChan {
            outputChan <- msg
        }
    }()
    
    return outputChan
}

func ConvertToAMQPDeliveryChan(inputChan <-chan interface{}) <-chan amqp.Delivery {
    outputChan := make(chan amqp.Delivery)

    go func() {
        defer close(outputChan)
        for msg := range inputChan {
            if amqpMsg, ok := msg.(amqp.Delivery); ok {
                outputChan <- amqpMsg
            } else {
                // Tratar o erro, se necessário
            }
        }
    }()

    return outputChan
}