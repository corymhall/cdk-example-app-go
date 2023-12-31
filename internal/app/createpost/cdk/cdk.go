package createpost

import (
	"log"
	"path"
	"runtime"

	c "github.com/corymhall/cdk-example-app-go/internal/pkg/cdk/constructs"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscodedeploy"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awsapigatewayv2 "github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewCreatePostFunction(scope constructs.Construct, id string, db awsdynamodb.ITable, api c.Api, monitor c.IMonitor) awslambda.IFunction {
	s := constructs.NewConstruct(scope, &id)

	env := &map[string]*string{
		"REGION":          awscdk.Stack_Of(scope).Region(),
		"POST_TABLE_NAME": db.TableName(),
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("error getting filename")
	}

	filepath := path.Join(path.Dir(filename), "../../../../cmd/createpost")

	handler := awslambdago.NewGoFunction(s, jsii.String("Handler"), &awslambdago.GoFunctionProps{
		Entry:       &filepath,
		Tracing:     awslambda.Tracing_ACTIVE,
		Environment: env,
		Bundling: &awslambdago.BundlingOptions{
			GoBuildFlags: &[]*string{
				jsii.String("-ldflags '-w -s -extldflags \"static\"'"),
				jsii.String("-a"),
			},
		},
		MemorySize: jsii.Number(512),
	})

	//
	// -------------------------------------------------------------------------
	// -----------------------add our route to our API Gateway -----------------
	// -------------------------------------------------------------------------
	//
	api.AddLambdaRoute("createPost", handler, c.AddRouteOptions{
		Path: jsii.String("/post"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_POST,
		},
	})

	//
	// -------------------------------------------------------------------------
	// -----------------------Create our deployment ----------------------------
	// -------------------------------------------------------------------------
	//
	alias := awslambda.NewAlias(s, jsii.String("Alias"), &awslambda.AliasProps{
		AliasName: jsii.String("live"),
		Version:   handler.CurrentVersion(),
	})

	awscodedeploy.NewLambdaDeploymentGroup(s, jsii.String("Canary"), &awscodedeploy.LambdaDeploymentGroupProps{
		Alias:            alias,
		DeploymentConfig: awscodedeploy.LambdaDeploymentConfig_CANARY_10PERCENT_10MINUTES(),
	})

	//
	// -------------------------------------------------------------------------
	// -----------------------Create our monitoring ----------------------------
	// -------------------------------------------------------------------------
	//

	// add our standard lambda metrics to our dashboard
	monitor.MonitorLambdaFunction(handler)

	// add a custom metric to track integration latecy
	metric := api.MetricIntegrationLatency(&awscloudwatch.MetricOptions{
		DimensionsMap: &map[string]*string{
			"ApiName":  api.HttpApiName(),
			"Stage":    api.DefaultStage().StageName(),
			"Method":   jsii.String("POST"),
			"Resource": jsii.String("/post"),
		},
		Period: awscdk.Duration_Minutes(jsii.Number(1)),
	})
	monitor.AddWidgets(
		awscloudwatch.NewGraphWidget(&awscloudwatch.GraphWidgetProps{
			Title:   jsii.String("POST /post (1-minute periods)"),
			Width:   jsii.Number(12),
			Stacked: jsii.Bool(false),
			Left:    &[]awscloudwatch.IMetric{metric},
		}),
	)

	//
	// -------------------------------------------------------------------------
	// -----------------------Grant access to DynamoDB -------------------------
	// -------------------------------------------------------------------------
	//

	// grant write access to the DynamoDB Table
	db.GrantWriteData(handler.GrantPrincipal())

	return handler
}
