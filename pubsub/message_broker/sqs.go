package message_broker

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func (s *SQS) SendMessage(message interface{}, queueName string) error {
	messageStr, ok := message.(string)
	if !ok {
		return fmt.Errorf("message is not a string")
	}

	params := &sqs.SendMessageInput{
		MessageBody: aws.String(messageStr),
		QueueUrl:    aws.String(queueName),
	}

	_, err := s.svc.SendMessage(params)
	return err
}

func (s *SQS) ReceiveMessage(queueName string) (<-chan interface{}, error) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueName),
		WaitTimeSeconds: aws.Int64(20),
	}

	resp, err := s.svc.ReceiveMessage(params)
	if err != nil {
		return nil, err
	}

	outputChan := make(chan interface{}, len(resp.Messages))

	go func() {
		defer close(outputChan)
		for _, msg := range resp.Messages {
			outputChan <- *msg.Body

			deleteParams := &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueName),
				ReceiptHandle: msg.ReceiptHandle,
			}

			_, err := s.svc.DeleteMessage(deleteParams)
			if err != nil {
				fmt.Println("Delete Error", err)
				continue
			}
		}
	}()

	return outputChan, nil
}

func (s *SQS) Close() {
	// Do nothing
}



