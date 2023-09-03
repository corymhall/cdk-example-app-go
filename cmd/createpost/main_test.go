package main

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/corymhall/cdk-example-app-go/internal/app/createpost/store"
	"github.com/corymhall/cdk-example-app-go/internal/pkg/models"
)

type mockPutItemAPI func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)

func (m mockPutItemAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m(ctx, params, optFns...)
}

func TestHandler(t *testing.T) {
	post := models.Post{
		Title:      "This is my post title",
		UserID:     "12345",
		Summary:    "This is my post summary",
		Content:    "this is a bunch of post content",
		PostStatus: "published",
		Categories: []string{"blog", "travel"},
		CreatedAt:  "2021-06-10T13:40:04.736Z",
		HeroImage:  "https://example.com/imageurl",
	}

	in, err := json.Marshal(post)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name    string
		success bool
		client  func(t *testing.T) store.DynamoDBPutItemAPI
		event   events.APIGatewayV2HTTPRequest
	}{
		{
			name:    "create post success",
			success: true,
			client: func(t *testing.T) store.DynamoDBPutItemAPI {
				return mockPutItemAPI(func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
					t.Helper()
					return &dynamodb.PutItemOutput{
						Attributes: params.Item,
					}, nil
				})
			},
			event: events.APIGatewayV2HTTPRequest{
				Body: string(json.RawMessage(in)),
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
