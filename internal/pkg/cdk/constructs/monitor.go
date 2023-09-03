package constructs

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudwatch"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awsapigatewayv2 "github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type MonitorProps struct {
	DashboardName string
}

type IMonitor interface {
	awscloudwatch.Dashboard
	MonitorLambdaFunction(function awslambda.IFunction)
	MonitorDynamoDBTable(table awsdynamodb.ITable)
	MonitorHttpApi(api awsapigatewayv2.HttpApi)
}

type Monitor struct {
	awscloudwatch.Dashboard
}

func NewMonitor(scope constructs.Construct, id string, props MonitorProps) IMonitor {
	s := constructs.NewConstruct(scope, &id)

	dash := awscloudwatch.NewDashboard(s, jsii.String("Dashboard"), &awscloudwatch.DashboardProps{
		DashboardName: &props.DashboardName,
	})

	return &Monitor{
		dash,
	}

}

func (m *Monitor) MonitorLambdaFunction(function awslambda.IFunction) {
	// invocations
	invocationMetric := function.MetricInvocations(&awscloudwatch.MetricOptions{
		DimensionsMap: &map[string]*string{
			"FunctionName": function.FunctionName(),
		},
		Statistic: jsii.String("sum"),
		Period:    awscdk.Duration_Minutes(jsii.Number(5)),
	})

	// errors
	errorsMetric := function.MetricErrors(&awscloudwatch.MetricOptions{
		DimensionsMap: &map[string]*string{
			"FunctionName": function.FunctionName(),
		},
		Statistic: jsii.String("sum"),
		Period:    awscdk.Duration_Minutes(jsii.Number(5)),
	})

	m.AddWidgets(
		awscloudwatch.NewGraphWidget(&awscloudwatch.GraphWidgetProps{
			Title: jsii.String("Errors/5min"),
			Width: jsii.Number(6),
			Left:  &[]awscloudwatch.IMetric{errorsMetric},
		}),
	)
	m.AddWidgets(
		awscloudwatch.NewGraphWidget(&awscloudwatch.GraphWidgetProps{
			Title: jsii.String("Invocations/5min"),
			Width: jsii.Number(6),
			Left:  &[]awscloudwatch.IMetric{invocationMetric},
		}),
	)
}

func (m *Monitor) MonitorHttpApi(api awsapigatewayv2.HttpApi) {
	countMetric := awscloudwatch.NewMetric(&awscloudwatch.MetricProps{
		DimensionsMap: &map[string]*string{
			"ApiName": api.HttpApiName(),
			"Stage":   api.DefaultStage().StageName(),
		},
		Namespace:  jsii.String("AWS/ApiGateway"),
		Period:     awscdk.Duration_Minutes(jsii.Number(1)),
		MetricName: jsii.String("Count"),
		Label:      jsii.String("Calls"),
		Color:      jsii.String("#1f77b4"),
		Statistic:  jsii.String("sum"),
	})

	m.AddWidgets(
		awscloudwatch.NewGraphWidget(&awscloudwatch.GraphWidgetProps{
			Title:   jsii.String("Overall Calls/1min"),
			Width:   jsii.Number(12),
			Stacked: jsii.Bool(false),
			Left:    &[]awscloudwatch.IMetric{countMetric},
		}),
	)

}

func (m *Monitor) MonitorDynamoDBTable(table awsdynamodb.ITable) {
	readMetric := awscloudwatch.NewMetric(&awscloudwatch.MetricProps{
		DimensionsMap: &map[string]*string{
			"TableName": table.TableName(),
		},
		Statistic:  jsii.String("sum"),
		Period:     awscdk.Duration_Minutes(jsii.Number(1)),
		Namespace:  jsii.String("AWS/DynamoDB"),
		MetricName: jsii.String("ConsumedReadCapacityUnits"),
		Label:      jsii.String("Consumed (Read)"),
	})

	writeMetric := awscloudwatch.NewMetric(&awscloudwatch.MetricProps{
		DimensionsMap: &map[string]*string{
			"TableName": table.TableName(),
		},
		Statistic:  jsii.String("sum"),
		Period:     awscdk.Duration_Minutes(jsii.Number(1)),
		Namespace:  jsii.String("AWS/DynamoDB"),
		MetricName: jsii.String("ConsumedWriteCapacityUnits"),
		Label:      jsii.String("Consumed (Write)"),
	})

	m.AddWidgets(
		awscloudwatch.NewGraphWidget(&awscloudwatch.GraphWidgetProps{
			Title:   jsii.String("Read Capacity Units/1min"),
			Width:   jsii.Number(12),
			Stacked: jsii.Bool(true),
			Left:    &[]awscloudwatch.IMetric{readMetric},
		}),
		awscloudwatch.NewGraphWidget(&awscloudwatch.GraphWidgetProps{
			Title:   jsii.String("Write Capacity Units/1min"),
			Width:   jsii.Number(12),
			Stacked: jsii.Bool(true),
			Left:    &[]awscloudwatch.IMetric{writeMetric},
		}),
	)
}
