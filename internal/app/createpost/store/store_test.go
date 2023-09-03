package store

import (
	"context"
	"strconv"
	"testing"

	"github.com/corymhall/cdk-example-app-go/internal/pkg/models"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type mockPutItemAPI func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)

func (m mockPutItemAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m(ctx, params, optFns...)
}

func TestCreatePost(t *testing.T) {
	cases := []struct {
		client    func(t *testing.T) DynamoDBPutItemAPI
		tableName string
		post      models.Post
		success   bool
	}{
		{
			client: func(t *testing.T) DynamoDBPutItemAPI {
				return mockPutItemAPI(func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
					t.Helper()
					if params.TableName == nil {
						t.Fatal("expect table name to not be nil")
					}
					if e, a := "testTable", *params.TableName; e != a {
						t.Errorf("expect %v, got %v", e, a)
					}
					return &dynamodb.PutItemOutput{
						Attributes: params.Item,
					}, nil
				})
			},
			tableName: "testTable",
			post: models.Post{
				Title:      "This is my post title",
				UserID:     "12345",
				Summary:    "This is my post summary",
				Content:    "this is a bunch of post content",
				PostStatus: "published",
				Categories: []string{"blog", "travel"},
				CreatedAt:  "2021-06-10T13:40:04.736Z",
				HeroImage:  "https://example.com/imageurl",
			},
			success: true,
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			err := CreatePost(ctx, tt.client(t), tt.tableName, tt.post)
			if (err == nil) != tt.success {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}
