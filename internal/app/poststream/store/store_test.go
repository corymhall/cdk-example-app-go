package store

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/corymhall/cdk-example-app-go/internal/pkg/models"
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

func TestGetPost(t *testing.T) {
	cases := []struct {
		client    func(t *testing.T) DynamoDBAPI
		tableName string
		pk        string
		success   bool
		expected  *models.Post
	}{
		{
			client: func(t *testing.T) DynamoDBAPI {
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
			tableName: "testTable",
			expected: &models.Post{
				PK:         "DETAILS",
				TotalPosts: 2,
			},
			pk:      "DETAILS",
			success: true,
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			res, err := GetPost(ctx, tt.client(t), tt.pk, tt.tableName)
			if !reflect.DeepEqual(res, tt.expected) {
				t.Fatalf("expect %v, got %v", *tt.expected, *res)
			}
			if (err == nil) != tt.success {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}

func TestUpdatePost(t *testing.T) {
	cases := []struct {
		client    func(t *testing.T) DynamoDBAPI
		tableName string
		post      models.Post
		success   bool
	}{
		{
			client: func(t *testing.T) DynamoDBAPI {
				return mockDynamoDBAPI{
					func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
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
					},
					func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
						return &dynamodb.GetItemOutput{}, nil
					},
				}
			},
			tableName: "testTable",
			post: models.Post{
				PK:         "DETAILS",
				TotalPosts: 2,
			},
			success: true,
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			err := UpdatePost(ctx, tt.client(t), tt.tableName, tt.post)
			if (err == nil) != tt.success {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}

func TestUpdatePostDetails(t *testing.T) {
	cases := []struct {
		client    func(t *testing.T) DynamoDBAPI
		tableName string
		action    int
		success   bool
	}{
		{
			client: func(t *testing.T) DynamoDBAPI {
				return mockDynamoDBAPI{
					func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						t.Helper()
						post := &models.Post{}
						if err := attributevalue.UnmarshalMap(params.Item, post); err != nil {
							t.Fatal(err)
						}

						if post.TotalPosts != 2 {
							t.Errorf("expect %v, got %v", post.TotalPosts, 2)
						}
						if params.TableName == nil {
							t.Fatal("expect table name to not be nil")
						}
						if e, a := "testTable", *params.TableName; e != a {
							t.Errorf("expect %v, got %v", e, a)
						}

						return &dynamodb.PutItemOutput{
							Attributes: map[string]types.AttributeValue{
								"pk": &types.AttributeValueMemberS{
									Value: "DETAILS",
								},
								"totalPosts": &types.AttributeValueMemberN{
									Value: "2",
								},
							},
						}, nil
					},
					func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
						return &dynamodb.GetItemOutput{
							Item: map[string]types.AttributeValue{
								"pk": &types.AttributeValueMemberS{
									Value: "DETAILS",
								},
								"totalPosts": &types.AttributeValueMemberN{
									Value: "1",
								},
							},
						}, nil
					},
				}
			},
			tableName: "testTable",
			action:    1,
			success:   true,
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			err := UpdatePostDetails(ctx, tt.client(t), tt.action, tt.tableName)
			if (err == nil) != tt.success {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}
