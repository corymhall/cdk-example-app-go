package store

import (
	"context"

	"github.com/corymhall/cdk-example-app-go/internal/pkg/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

func UpdatePostDetails(ctx context.Context, api DynamoDBAPI, action int, tableName string) error {
	old, err := GetPost(ctx, api, "DETAILS", tableName)
	if err != nil {
		return err
	}

	if old.PK == "" {
		old.PK = "DETAILS"
	}

	old.TotalPosts = old.TotalPosts + action

	if err := UpdatePost(ctx, api, tableName, *old); err != nil {
		return err
	}

	return nil
}

func UpdatePost(ctx context.Context, api DynamoDBAPI, tableName string, post models.Post) error {
	val, err := attributevalue.MarshalMap(post)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:                   val,
		TableName:              aws.String(tableName),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}

	_, err = api.PutItem(ctx, input)
	if err != nil {
		return errors.Wrap(err, "error updating item in dynamodb")
	}

	return nil
}

func GetPost(ctx context.Context, api DynamoDBAPI, pk, tableName string) (*models.Post, error) {
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
