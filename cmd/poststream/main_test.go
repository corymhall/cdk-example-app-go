package main

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/corymhall/cdk-example-app-go/internal/app/poststream/store"
)

type mockDynamoDBAPI struct {
	mockPutItemAPI func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	mockGetItemAPI func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

func (m mockDynamoDBAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.mockPutItemAPI(ctx, params, optFns...)
}

func (m mockDynamoDBAPI) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return m.mockGetItemAPI(ctx, params, optFns...)
}

func TestHandler(t *testing.T) {
	pki := []byte(`{ "S": "post-title-1"}`)
	ski := []byte(`{ "S": "published#private"}`)
	var pk events.DynamoDBAttributeValue
	var sk events.DynamoDBAttributeValue
	_ = json.Unmarshal(pki, &pk)
	_ = json.Unmarshal(ski, &sk)

	cases := []struct {
		name    string
		success bool
		client  func(t *testing.T) store.DynamoDBAPI
		event   events.DynamoDBEvent
	}{
		{
			name:    "create post success",
			success: true,
			client: func(t *testing.T) store.DynamoDBAPI {
				return mockDynamoDBAPI{
					func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						t.Helper()
						return &dynamodb.PutItemOutput{
							Attributes: params.Item,
						}, nil
					},
					func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
						t.Helper()
						return &dynamodb.GetItemOutput{
							Item: map[string]types.AttributeValue{
								"pk": &types.AttributeValueMemberS{
									Value: "DETAILS",
								},
								"totalPosts": &types.AttributeValueMemberN{
									Value: "2",
								},
							},
						}, nil
					},
				}
			},
			event: events.DynamoDBEvent{
				Records: []events.DynamoDBEventRecord{
					{
						EventName: "INSERT",
						Change: events.DynamoDBStreamRecord{
							Keys: map[string]events.DynamoDBAttributeValue{
								"pk": pk,
								"sk": sk,
							},
						},
					},
				},
			},
		},
	}
	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			st := &Store{
				client:    tt.client(t),
				tableName: "TestTable",
			}
			err := st.Handler(ctx, tt.event)
			if (err == nil) != tt.success {
				t.Fatalf("expect no error, got %v", err)
			}

		})
	}

}
