package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/corymhall/cdk-example-app-go/internal/app/poststream/store"
	"github.com/kelseyhightower/envconfig"
)

var client store.DynamoDBAPI

type Spec struct {
	Region    string `default:"us-east-2"`
	TableName string `required:"false" default:"TESTING" split_words:"true"`
}

type Store struct {
	client    store.DynamoDBAPI
	tableName string
}

func (s *Store) Handler(ctx context.Context, event events.DynamoDBEvent) error {

	for _, record := range event.Records {
		if record.EventName != "MODIFY" {
			pk := record.Change.Keys["pk"].String()

			// don't process changes to the DETAILS record
			if pk == "DETAILS" {
				return nil
			}

			switch record.EventName {
			case "INSERT":
				if err := store.UpdatePostDetails(ctx, s.client, 1, s.tableName); err != nil {
					return err
				}
			case "REMOVE":
				if err := store.UpdatePostDetails(ctx, s.client, -1, s.tableName); err != nil {
					return err
				}
			}
		}
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
