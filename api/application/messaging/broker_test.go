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
func TestSQSMessageBroker_SendMessagesInParallel(t *testing.T) {

	ctx := context.TODO()
	container := testcontainer.SportLinkContainer(t, ctx)
	defer container.Terminate(ctx)
	sqsClient := testcontainer.GetSqsClient(t, container, ctx)
	broker := messaging.NewBroker(sqsClient, fmt.Sprintf("http://localhost:4566/000000000000/%s", sqsQueueName))

	testCases := []struct {
		name                 string
		batch                []messaging.Message
		expectedSentMessages int
		numIterations        int
		assertions           func(t *testing.T, err error)
	}{
		{
			name: "send a batch of messages in parallel with the max amount of messages successfully 2 times",
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
			numIterations:        2,
			expectedSentMessages: 2 * 10,
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "send a batch of messages in parallel successfully 20 times",
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
			numIterations:        20,
			expectedSentMessages: 4 * 20,
			assertions: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			errChan := make(chan error, testCase.numIterations)
			outputChan := make(chan messaging.SendMessagesOutput, testCase.numIterations)

			// when
			for i := 0; i < testCase.numIterations; i++ {
				go func() {
					// when
					output, err := broker.SendMessages(ctx, testCase.batch)
					if err != nil {
						errChan <- err
					} else {
						outputChan <- output
					}
				}()
			}

			// then
			totalFailed, totalSucceeded := collectResultFromChannel(t, testCase.numIterations, errChan, outputChan)
			assert.Equal(t, testCase.expectedSentMessages, totalSucceeded)
			assert.Equal(t, 0, totalFailed, "Expected 0 failed messages")
		})
	}
}

func collectResultFromChannel(t *testing.T, numIterations int, errChan chan error, outputChan chan messaging.SendMessagesOutput) (int, int) {
	totalSucceeded := 0
	totalFailed := 0
	for i := 0; i < numIterations; i++ {
		select {
		case err := <-errChan:
			if err != nil {
				t.Errorf("Error sending messages: %v", err)
				totalFailed++
			}
		case output := <-outputChan:
			totalSucceeded += output.Succeeded
			totalFailed += output.Failed
		}
	}
	return totalFailed, totalSucceeded
}
