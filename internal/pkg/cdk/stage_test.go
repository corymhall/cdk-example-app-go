package cdk

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
	"github.com/stretchr/testify/assert"
)

var stage awscdk.Stage

func TestMain(m *testing.M) {
	app := awscdk.NewApp(nil)
	stage = NewStage(app, "CreatePost", awscdk.StageProps{})

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestAllStacksCreated(t *testing.T) {
	stacks := stage.Synth(nil).StacksRecursively()
	assert.Equal(t, len(*stacks), 3)
}

func TestDynamoDB(t *testing.T) {
	bytes, err := json.Marshal(stage.Synth(nil).GetStackByName(jsii.String("CreatePost-DatastoreStack")).Template())
	if err != nil {
		t.Error(err)
	}
	template := assertions.Template_FromString(jsii.String(string(bytes)), &assertions.TemplateParsingOptions{})
	resources := template.FindResources(jsii.String("AWS::DynamoDB::Table"), nil)
	if _, ok := (*resources)["PostDatastore470A9474"]; !ok {
		t.Error("DynamoDB table PostDatastore470A9474 does not exist")
	}
}

func TestApi(t *testing.T) {
	bytes, err := json.Marshal(stage.Synth(nil).GetStackByName(jsii.String("CreatePost-APIStack")).Template())
	if err != nil {
		t.Error(err)
	}

	template := assertions.Template_FromString(jsii.String(string(bytes)), &assertions.TemplateParsingOptions{})
	resources := template.FindResources(jsii.String("AWS::ApiGatewayV2::Api"), nil)
	if _, ok := (*resources)["ApiF70053CD"]; !ok {
		t.Error("Api Gateway API ApiF70053CD does not exist")
	}
}
