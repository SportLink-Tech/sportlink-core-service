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
	SendMessages(ctx context.Context, batch []BatchMessage) error
	ReceiveMessages(ctx context.Context, batchSize int) ([]string, error)
}

// BatchMessage represents a single message in a batch
type BatchMessage struct {
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

// SendMessages sends a pre-formed batch of messages to the SQS queue
func (broker *SQSMessageBroker) SendMessages(ctx context.Context, batch []BatchMessage) error {
	if len(batch) > 10 {
		return fmt.Errorf("batch size exceeds SQS limit of 10 messages per batch")
	}

	entries := make([]types.SendMessageBatchRequestEntry, len(batch))
	for i, message := range batch {
		entries[i] = types.SendMessageBatchRequestEntry{
			Id:          aws.String(message.Id),
			MessageBody: aws.String(message.Message),
		}
	}

	msg := &sqs.SendMessageBatchInput{
		QueueUrl: aws.String(broker.queueUrl),
		Entries:  entries,
	}

	result, err := broker.client.SendMessageBatch(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to send batch of messages: %w", err)
	}

	if len(result.Failed) > 0 {
		for _, failure := range result.Failed {
			log.Printf("Failed to send message ID: %s, Reason: %s", *failure.Id, *failure.Message)
		}
		return fmt.Errorf("some messages failed to send")
	}

	return nil
}

// ReceiveMessages receives a batch of messages from the SQS queue
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
