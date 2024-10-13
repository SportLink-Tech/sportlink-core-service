package messaging_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sportlink/api/application/messaging"
	"sportlink/dev/testcontainer"
	"testing"
)

const sqsQueueName = "sportlink-news"

func TestSQSMessageBroker_SendMessage(t *testing.T) {

	ctx := context.TODO()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	sqsClient := testcontainer.GetSqsClient(t, container, ctx)
	broker := messaging.NewBroker(sqsClient, fmt.Sprintf("http://localhost:4566/000000000000/%s", sqsQueueName))

	testCases := []struct {
		name       string
		message    string
		assertions func(t *testing.T, err error)
	}{
		{
			name:    "send a message successfully",
			message: `{"description":"Sport link core service"}`,
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "send an empty message must retrieve an error",
			assertions: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "message is empty")
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// when
			err := broker.SendMessage(ctx, testCase.message)

			// then
			testCase.assertions(t, err)
		})
	}
}

// Use the following command in order to see the amount of messages:
// ```shell
// awslocal sqs get-queue-attributes --queue-url=http://localhost:4566/000000000000/sportlink-news --attribute-names ApproximateNumberOfMessages
// ```
func TestSQSMessageBroker_SendMessages(t *testing.T) {

	ctx := context.TODO()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	sqsClient := testcontainer.GetSqsClient(t, container, ctx)
	broker := messaging.NewBroker(sqsClient, fmt.Sprintf("http://localhost:4566/000000000000/%s", sqsQueueName))

	testCases := []struct {
		name       string
		batch      []messaging.Message
		assertions func(t *testing.T, err error)
	}{
		{
			name: "send a batch of messages successfully",
			batch: []messaging.Message{
				{
					Id:      "1",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "2",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "3",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "4",
					Message: `{"description":"Sport link core service"}`,
				},
			},
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "send a batch of messages with more than the max amount of messages available",
			batch: []messaging.Message{
				{
					Id:      "1",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "2",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "3",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "4",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "5",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "6",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "7",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "8",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "9",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "10",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "11",
					Message: `{"description":"Sport link core service"}`,
				},
			},
			assertions: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "batch size exceeds SQS limit of 10 messages per batch")
			},
		},
		{
			name: "send a batch of messages with the max amount of messages successfully",
			batch: []messaging.Message{
				{
					Id:      "1",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "2",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "3",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "4",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "5",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "6",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "7",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "8",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "9",
					Message: `{"description":"Sport link core service"}`,
				},
				{
					Id:      "10",
					Message: `{"description":"Sport link core service"}`,
				},
			},
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "send a batch of messages with a single message successfully",
			batch: []messaging.Message{
				{
					Id:      "1",
					Message: `{"description":"Sport link core service"}`,
				},
			},
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// when
			err := broker.SendMessages(ctx, testCase.batch)

			// then
			testCase.assertions(t, err)
		})
	}
}
