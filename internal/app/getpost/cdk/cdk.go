package getpost

import (
	"log"
	"path"
	"runtime"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscodedeploy"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awsapigatewayv2 "github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	integrations "github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	c "github.com/corymhall/cdk-example-app-go/internal/pkg/cdk/constructs"
)

func NewGetPostFunction(scope constructs.Construct, id string, db awsdynamodb.ITable, api awsapigatewayv2.HttpApi, monitor c.IMonitor) awslambda.IFunction {
	s := constructs.NewConstruct(scope, &id)

	env := &map[string]*string{
		"REGION":          awscdk.Stack_Of(scope).Region(),
		"POST_TABLE_NAME": db.TableName(),
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("error getting filename")
	}

	filepath := path.Join(path.Dir(filename), "../../../../cmd/getpost")

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
	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: integrations.NewHttpLambdaIntegration(jsii.String("getPost"), handler, &integrations.HttpLambdaIntegrationProps{}),
		Path:        jsii.String("/post/{postId}"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
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
			"Method":   jsii.String("GET"),
			"Resource": jsii.String("/post/{postId}"),
		},
		Period: awscdk.Duration_Minutes(jsii.Number(1)),
	})
	monitor.AddWidgets(
		awscloudwatch.NewGraphWidget(&awscloudwatch.GraphWidgetProps{
			Title:   jsii.String("GET /post (1-minute periods)"),
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

	// grant read access to the DynamoDB Table
	db.GrantReadData(handler.GrantPrincipal())

	return handler
}
