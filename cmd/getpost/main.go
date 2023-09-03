package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/corymhall/cdk-example-app-go/internal/app/getpost/store"
	"github.com/corymhall/cdk-example-app-go/internal/pkg/models"
	"github.com/kelseyhightower/envconfig"
)

var client store.DynamoDBGetItemAPI

type Store struct {
	client    store.DynamoDBGetItemAPI
	tableName string
}

func (s *Store) Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*models.Post, error) {
	postId := event.PathParameters["postId"]

	res, err := store.GetPost(ctx, s.client, postId, s.tableName)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type Spec struct {
	Region    string `default:"us-east-2"`
	TableName string `required:"false" default:"TESTING" split_words:"true"`
}

func main() {
	var s Spec
	if err := envconfig.Process("post", &s); err != nil {
		log.Fatal(err.Error())
	}
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	client = dynamodb.NewFromConfig(cfg)
	st := &Store{
		client:    client,
		tableName: s.TableName,
	}

	lambda.Start(st.Handler)
}
