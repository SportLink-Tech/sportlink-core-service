package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDb[E any] struct {
	client         dynamodb.Client
	requestBuilder RequestBuilder[E]
	tableName      string
}

func NewDynamoDb[E any](client dynamodb.Client, tableName string) *DynamoDb[E] {
	requestBuilder := NewRequestBuilder[E]()
	return &DynamoDb[E]{
		client:         client,
		tableName:      tableName,
		requestBuilder: requestBuilder,
	}
}

func (k *DynamoDb[E]) Get(ctx context.Context, key map[string]types.AttributeValue) (*E, error) {
	r, err := k.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(k.tableName),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	if r.Item == nil {
		return nil, nil
	}

	var item E
	if err := attributevalue.UnmarshalMap(r.Item, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (k *DynamoDb[E]) Save(ctx context.Context, itemDto E) error {
	item, err := attributevalue.MarshalMap(itemDto)
	if err != nil {
		return err
	}

	_, err = k.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(k.tableName),
		Item:      item,
	})

	return err
}

func (k *DynamoDb[E]) SaveAll(ctx context.Context, items []E, batchSize int) (SaveAllOutput, error) {
	output := SaveAllOutput{UnprocessedItems: 0, ProcessedItems: 0}
	writeRequests, err := k.requestBuilder.WriteRequest(items)
	if err != nil {
		return output, err
	}

	for start := 0; start < len(writeRequests); start += batchSize {
		end := start + batchSize
		if end > len(writeRequests) {
			end = len(writeRequests)
		}

		batchWriteRequest := writeRequests[start:end]
		batchWriteItemOutput, err := k.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				k.tableName: batchWriteRequest,
			},
		})

		if err != nil {
			return output, err
		}

		output.UnprocessedItems += len(batchWriteItemOutput.UnprocessedItems)
		output.ProcessedItems += len(batchWriteRequest) + len(batchWriteItemOutput.UnprocessedItems)
	}

	return output, nil
}
