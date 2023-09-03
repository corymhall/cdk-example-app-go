package poststream

import (
	"log"
	"path"
	"runtime"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	c "github.com/corymhall/cdk-example-app-go/internal/pkg/cdk/constructs"
)

func NewPostStreamFunction(scope constructs.Construct, id string, db awsdynamodb.ITable, monitor c.IMonitor) awslambda.IFunction {
	s := constructs.NewConstruct(scope, &id)

	env := &map[string]*string{
		"REGION":          awscdk.Stack_Of(scope).Region(),
		"POST_TABLE_NAME": db.TableName(),
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("error getting filename")
	}

	filepath := path.Join(path.Dir(filename), "../../../../cmd/poststream")

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

	dlq := awssqs.NewQueue(s, jsii.String("DLQ"), &awssqs.QueueProps{})

	handler.AddEventSource(
		awslambdaeventsources.NewDynamoEventSource(db, &awslambdaeventsources.DynamoEventSourceProps{
			StartingPosition: awslambda.StartingPosition_LATEST,
			OnFailure:        awslambdaeventsources.NewSqsDlq(dlq),
		}),
	)

	db.GrantReadWriteData(handler.GrantPrincipal())
	db.GrantStreamRead(handler.GrantPrincipal())

	monitor.MonitorLambdaFunction(handler)

	return handler
}
