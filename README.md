# AWS CDK Golang Example App

This contains the example app for the blog post [Application and Infrastructure together with the AWS CDK for Go]()

## Deploying the example

If you would like to deploy this example application you can follow the below steps.


### Setup Environment

By default this application will use the `CDK_DEFAULT_ACCOUNT` and `CDK_DEFAULT_REGION` environment variables
for both the pipeline and application. If you would like to change this behavior, for example to deploy
the pipeline to one account and the application to another, you can edit the `Env` property for each stack
in [cdk-go-app.go](./cdk-go-app.go)

_For example_
```go
pipelineStack := awscdk.NewStack(app, jsii.String("PostStoreGo-DeliveryPipeline"), &awscdk.StackProps{
	Env: NewEnv("111111111111", "us-east-1"),
})
```


Next you will need to bootstrap the accounts you are deploying to. If you are working with accounts that
have already been bootstrapped then you can skip this step.

_note: if you are deploying to multiple accounts you will need to run the command separately for each account_
```bash
cdk bootstrap
```


### Deploy

#### Without pipeline

If you want to deploy the app without deploying through the pipeline you can target the stage

```bash
cdk deploy 'AppStage/*'
```

This will deploy the `AppStage/DatastoreStack`, `AppStage/APIStack`, & `AppStage/MonitorStack`.

#### With the pipeline

If you want to deploy the solution with the pipeline you'll need to first create your own GitHub repo
for this application.

The next thing you will need is a CodeStar Connection to GitHub. Follow the
instructions [here](https://docs.aws.amazon.com/codepipeline/latest/userguide/connections-github.html)

Once you have the repo you'll need to edit the pipeline stack to point to your repo.

[cdk-go-app.go](./cdk-go-app.go)
```go
pipeline := cdk.NewPipeline(pipelineStack, "PostStoreGo-DeliveryPipeline", cdk.PipelineProps{
	Name: "PostStoreGo",
	Owner: "UPDATE_WITH_OWNER",
	Repo: "UPDATE_WITH_REPO_NAME",
	ConnectionArn: "ARN_OF_CODESTAR_CONNECTION",
})
```

Make sure you commit all of you changes to git and then deploy the pipeline stack
to kickoff the pipeline.

```bash
cdk deploy 'PostStoreGo-DeliveryPipeline'
```

### Test Application

If you are deploying the application locally you should see an `Outputs` list that contains
something similar to `AppStageAPIStack81234557.APIURL`. The value is the URL of your API Gateway.

#### Create a new post

```bash
curl --request POST \
  --url REPLACE_WITH_APIURL/post \
  --header 'Content-Type: application/json' \
  --data '{
        "title": "Test Post",
        "pk": "test-post"
}'
```

#### Viewing Post

```bash
curl REPLACE_WITH_APIURL/post/test-post
```

