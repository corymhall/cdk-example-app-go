package store

import (
	"context"

	"github.com/corymhall/cdk-example-app-go/internal/pkg/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

type DynamoDBPutItemAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func CreatePost(ctx context.Context, api DynamoDBPutItemAPI, tableName string, post models.Post) error {
	cond := expression.Name("pk").AttributeNotExists()
	expr, err := expression.NewBuilder().WithCondition(cond).Build()
	if err != nil {
		return err
	}

	val, err := attributevalue.MarshalMap(post)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:                     val,
		TableName:                aws.String(tableName),
		ConditionExpression:      expr.Condition(),
		ExpressionAttributeNames: expr.Names(),
		ReturnConsumedCapacity:   types.ReturnConsumedCapacityTotal,
	}

	_, err = api.PutItem(ctx, input)
	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if errors.As(err, &condErr) {
			return errors.Wrapf(err, "%s: error creating item", "")
		}
		return errors.Wrap(err, "error creating item in dynamodb")
	}

	return nil
}
