package cdk

import (
	createpost "github.com/corymhall/cdk-example-app-go/internal/app/createpost/cdk"
	getpost "github.com/corymhall/cdk-example-app-go/internal/app/getpost/cdk"
	poststream "github.com/corymhall/cdk-example-app-go/internal/app/poststream/cdk"
	c "github.com/corymhall/cdk-example-app-go/internal/pkg/cdk/constructs"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	awsapigatewayv2 "github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewStage(scope constructs.Construct, id string, props awscdk.StageProps) awscdk.Stage {
	stage := awscdk.NewStage(scope, &id, &props)

	//---------------------//
	//----Monitor Stack----//
	//---------------------//

	monitorStack := awscdk.NewStack(stage, jsii.String("MonitorStack"), &awscdk.StackProps{})
	monitor := c.NewMonitor(monitorStack, "Monitor", c.MonitorProps{
		DashboardName: "PostStore-Dashboard",
	})

	//-----------------------//
	//----Datastore Stack----//
	//-----------------------//

	dynamodbStack := awscdk.NewStack(stage, jsii.String("DatastoreStack"), &awscdk.StackProps{
		TerminationProtection: jsii.Bool(true),
	})

	table := awsdynamodb.NewTable(dynamodbStack, jsii.String("PostDatastore"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("pk"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		Encryption:  awsdynamodb.TableEncryption_AWS_MANAGED,
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
		Stream:      awsdynamodb.StreamViewType_NEW_AND_OLD_IMAGES,
	})

	monitor.MonitorDynamoDBTable(table)

	//-----------------------//
	//-------Api Stack-------//
	//-----------------------//

	apiStack := awscdk.NewStack(stage, jsii.String("APIStack"), &awscdk.StackProps{})

	api := awsapigatewayv2.NewHttpApi(apiStack, jsii.String("Api"), &awsapigatewayv2.HttpApiProps{})
	monitor.MonitorHttpApi(api)

	awscdk.NewCfnOutput(apiStack, jsii.String("API-URL"), &awscdk.CfnOutputProps{
		Value: api.Url(),
	})

	createpost.NewCreatePostFunction(apiStack, "PostApi-CreatePost", table, api, monitor)

	poststream.NewPostStreamFunction(apiStack, "PostApi-PostStream", table, monitor)

	getpost.NewGetPostFunction(apiStack, "PostApi-GetPost", table, api, monitor)

	return stage
}
