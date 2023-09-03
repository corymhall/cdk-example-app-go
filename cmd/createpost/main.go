package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/corymhall/cdk-example-app-go/internal/app/createpost/store"
	"github.com/corymhall/cdk-example-app-go/internal/pkg/models"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/kelseyhightower/envconfig"
)

var client store.DynamoDBPutItemAPI

type Spec struct {
	Region    string `default:"us-east-2"`
	TableName string `required:"false" default:"TESTING" split_words:"true"`
}

type Store struct {
	client    store.DynamoDBPutItemAPI
	tableName string
}

func (s *Store) Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) error {
	post := &models.Post{}
	if err := json.Unmarshal([]byte(event.Body), post); err != nil {
		return err
	}

	if err := store.CreatePost(ctx, s.client, s.tableName, *post); err != nil {
		return err
	}

	return nil
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
