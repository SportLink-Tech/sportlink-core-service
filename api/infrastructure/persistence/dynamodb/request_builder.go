package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type RequestBuilder[E any] interface {
	WriteRequest(items []E) ([]types.WriteRequest, error)
}

type DynamoRequestBuilder[E any] struct {
}

func NewRequestBuilder[E any]() *DynamoRequestBuilder[E] {
	return &DynamoRequestBuilder[E]{}
}

func (r *DynamoRequestBuilder[E]) WriteRequest(items []E) ([]types.WriteRequest, error) {
	var writeRequests []types.WriteRequest

	for _, itemDto := range items {
		item, err := attributevalue.MarshalMap(itemDto)
		if err != nil {
			return nil, err
		}

		writeRequest := types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: item,
			},
		}

		writeRequests = append(writeRequests, writeRequest)
	}
	return writeRequests, nil
}
