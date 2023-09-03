package store

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/corymhall/cdk-example-app-go/internal/pkg/models"
	"github.com/pkg/errors"
)

type DynamoDBGetItemAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

func GetPost(ctx context.Context, api DynamoDBGetItemAPI, pk, tableName string) (*models.Post, error) {
	post := &models.Post{}
	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{
				Value: pk,
			},
		},
		TableName:              aws.String(tableName),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}

	res, err := api.GetItem(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "error getting post")
	}

	if err := attributevalue.UnmarshalMap(res.Item, post); err != nil {
		return nil, err
	}

	return post, nil
}
