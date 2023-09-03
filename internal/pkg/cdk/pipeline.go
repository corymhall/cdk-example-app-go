package cdk

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	"github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type PipelineProps struct {
	awscdk.StackProps
	Name          string
	Owner         string
	Repo          string
	ConnectionArn string
}

func NewPipeline(scope constructs.Construct, id string, props PipelineProps) pipelines.CodePipeline {
	s := constructs.NewConstruct(scope, &id)

	pipeline := pipelines.NewCodePipeline(s, jsii.String("Pipeline"), &pipelines.CodePipelineProps{
		PipelineName:     &props.Name,
		CrossAccountKeys: jsii.Bool(true),
		Synth: pipelines.NewShellStep(jsii.String("Synth"), &pipelines.ShellStepProps{
			Input: pipelines.CodePipelineSource_Connection(jsii.String(fmt.Sprintf("%s/%s", props.Owner, props.Repo)), jsii.String("main"), &pipelines.ConnectionSourceOptions{
				ConnectionArn: &props.ConnectionArn,
			}),
			Commands: jsii.Strings(
				"cdk synth -v",
			),
			InstallCommands: jsii.Strings(
				"npm install -g aws-cdk",
			),
		}),
		CodeBuildDefaults: &pipelines.CodeBuildOptions{
			BuildEnvironment: &awscodebuild.BuildEnvironment{
				BuildImage: awscodebuild.LinuxBuildImage_STANDARD_7_0(),
			},
		},
	})

	return pipeline
}
