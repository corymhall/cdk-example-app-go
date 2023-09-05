package constructs

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ApiGatewayService struct {
	constructs.Construct
	service        awsecs.FargateService
	taskDefinition awsecs.IFargateTaskDefinition
}

func NewApiGatewayService(scope constructs.Construct, id string) ApiGatewayService {
	s := constructs.NewConstruct(scope, &id)
	api := ApiGatewayService{}

	logGroup := awslogs.NewLogGroup(s, jsii.String("LogGroup"), &awslogs.LogGroupProps{
		Retention: awslogs.RetentionDays_ONE_MONTH,
	})

	taskDef := awsecs.NewFargateTaskDefinition(s, jsii.String("TaskDef"), &awsecs.FargateTaskDefinitionProps{
		MemoryLimitMiB: jsii.Number(512),
	})

	taskDef.AddExtension(XRayExtension{
		Extension: Extension{
			logGroup: logGroup,
		},
	})

	return api
}

type Extension struct {
	logGroup awslogs.LogGroup
}

type XRayExtension struct {
	Extension
}

func (e XRayExtension) Extend(taskDef awsecs.TaskDefinition) {
	taskDef.AddContainer(jsii.String("xray"), &awsecs.ContainerDefinitionOptions{
		Image:                awsecs.AssetImage_FromRegistry(jsii.String("amazon/aws-xray-daemon:latest"), &awsecs.RepositoryImageProps{}),
		Essential:            jsii.Bool(true),
		MemoryReservationMiB: jsii.Number(256),
		PortMappings: &[]*awsecs.PortMapping{
			{
				ContainerPort: jsii.Number(2000),
				Protocol:      awsecs.Protocol_UDP,
			},
		},
		Environment: &map[string]*string{
			"AWS_REGION": awscdk.Stack_Of(taskDef).Region(),
		},
		HealthCheck: &awsecs.HealthCheck{
			Command: jsii.Strings(
				"CMD-SHELL",
				"curl -s http://localhost:2000",
			),
			StartPeriod: awscdk.Duration_Seconds(jsii.Number(10)),
			Interval:    awscdk.Duration_Seconds(jsii.Number(5)),
			Timeout:     awscdk.Duration_Seconds(jsii.Number(2)),
			Retries:     jsii.Number(3),
		},
		Logging: awsecs.NewAwsLogDriver(&awsecs.AwsLogDriverProps{
			StreamPrefix: jsii.String("xray"),
			LogGroup:     e.logGroup,
		}),
	})
	taskDef.TaskRole().AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AWSXRayDaemonWriteAccess")))
}
