package main

import (
	"os"

	"github.com/corymhall/cdk-example-app-go/internal/pkg/cdk"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
	"github.com/aws/jsii-runtime-go"
)

func NewEnv(account, region *string) *awscdk.Environment {
	if account == nil {
		account = jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT"))
	}

	if region == nil {
		region = jsii.String(os.Getenv("CDK_DEFAULT_REGION"))
	}
	return &awscdk.Environment{
		Account: account,
		Region:  region,
	}
}

func main() {
	app := awscdk.NewApp(nil)

	pipelineStack := awscdk.NewStack(app, jsii.String("PostStoreGo-DeliveryPipeline"), &awscdk.StackProps{
		Env: NewEnv(nil, nil),
	})

	pipeline := cdk.NewPipeline(pipelineStack, "PostStoreGo-DeliveryPipeline", cdk.PipelineProps{
		Name: "PostStoreGo",
	})

	appStage := cdk.NewStage(app, "AppStage", awscdk.StageProps{
		Env: NewEnv(nil, nil),
	})

	awscdk.Tags_Of(appStage).Add(
		jsii.String("project"),
		jsii.String("blog-go"),
		&awscdk.TagProps{},
	)
	awscdk.Tags_Of(appStage).Add(
		jsii.String("application"),
		jsii.String("PostStoreGo"),
		&awscdk.TagProps{},
	)

	pipeline.AddStage(appStage, &pipelines.AddStageOpts{})

	app.Synth(nil)
}
