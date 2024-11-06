package messaging

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Broker interface with methods to send and receive messages
type Broker interface {
	SendMessage(ctx context.Context, message string) error
	SendMessages(ctx context.Context, batch []Message) (SendMessagesOutput, error)
	ReceiveMessages(ctx context.Context, batchSize int) ([]string, error)
}

type Message struct {
	Id      string
	Message string
}

type SQSMessageBroker struct {
	client   *sqs.Client
	queueUrl string
}

func NewBroker(client *sqs.Client, queueUrl string) Broker {
	return &SQSMessageBroker{
		client:   client,
		queueUrl: queueUrl,
	}
}

func (broker *SQSMessageBroker) SendMessage(ctx context.Context, message string) error {
	if message == "" {
		return fmt.Errorf("message is empty")
	}
	sendMsgInput := &sqs.SendMessageInput{
		QueueUrl:    &broker.queueUrl,
		MessageBody: &message,
	}

	_, err := broker.client.SendMessage(ctx, sendMsgInput)
	if err != nil {
		return err
	}

	return nil
}

func (broker *SQSMessageBroker) SendMessages(ctx context.Context, batch []Message) (SendMessagesOutput, error) {
	if len(batch) > 10 {
		return SendMessagesOutput{}, fmt.Errorf("batch size exceeds SQS limit of 10 messages per batch")
	}

	entries := make([]types.SendMessageBatchRequestEntry, len(batch))
	for i, message := range batch {
		entries[i] = types.SendMessageBatchRequestEntry{
			Id:          aws.String(message.Id),
			MessageBody: aws.String(message.Message),
		}
	}

	msg := &sqs.SendMessageBatchInput{
		QueueUrl: &broker.queueUrl,
		Entries:  entries,
	}

	result, err := broker.client.SendMessageBatch(ctx, msg)
	if err != nil {
		return SendMessagesOutput{}, fmt.Errorf("failed to send batch of messages: %w", err)
	}

	if len(result.Failed) > 0 {
		for _, failure := range result.Failed {
			log.Printf("Failed to send message Id: %s, Reason: %s", *failure.Id, *failure.Message)
		}
		return SendMessagesOutput{
			Succeeded: len(result.Successful),
			Failed:    len(result.Failed),
		}, fmt.Errorf("some messages failed to send")
	}

	return SendMessagesOutput{
		Succeeded: len(result.Successful),
		Failed:    len(result.Failed),
	}, nil
}

func (broker *SQSMessageBroker) ReceiveMessages(ctx context.Context, batchSize int) ([]string, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(broker.queueUrl),
		MaxNumberOfMessages: int32(batchSize),
	}

	result, err := broker.client.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, err
	}

	messages := make([]string, 0, len(result.Messages))
	for _, msg := range result.Messages {
		messages = append(messages, *msg.Body)
	}

	return messages, nil
}
