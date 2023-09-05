package constructs

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsservicediscovery"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type Api interface {
	constructs.Construct
	awsec2.IConnectable
	awscdkapigatewayv2alpha.HttpApi
	AddServiceRoute(string, awsservicediscovery.IService, AddRouteOptions) *[]awscdkapigatewayv2alpha.HttpRoute
	AddLambdaRoute(string, awslambda.IFunction, AddRouteOptions) *[]awscdkapigatewayv2alpha.HttpRoute
}

type Option func(*ServiceApi) error

func WithVpc(vpc awsec2.IVpc) func(*ServiceApi) {
	return func(a *ServiceApi) {
		sg := awsec2.NewSecurityGroup(*a, jsii.String("VpcLinkSG"), &awsec2.SecurityGroupProps{
			Vpc:         vpc,
			Description: jsii.String("Security group for VPC Link connections"),
		})
		vpcLink := awscdkapigatewayv2alpha.NewVpcLink(*a, jsii.String("VpcLink"), &awscdkapigatewayv2alpha.VpcLinkProps{
			Vpc:            vpc,
			SecurityGroups: &[]awsec2.ISecurityGroup{sg},
		})
		a.vpcLink = &vpcLink
		a.Connections().AddSecurityGroup(sg)
	}
}

type ServiceApi struct {
	awscdkapigatewayv2alpha.HttpApi
	vpcLink *awscdkapigatewayv2alpha.VpcLink
}

func NewApi(scope constructs.Construct, id string, options ...Option) Api {
	construct := constructs.NewConstruct(scope, &id)
	service := &ServiceApi{}
	for _, o := range options {
		o(service)
	}
	base := awscdkapigatewayv2alpha.NewHttpApi(construct, jsii.String("HttpApi"), &awscdkapigatewayv2alpha.HttpApiProps{
		CreateDefaultStage: jsii.Bool(true),
	})
	service.HttpApi = base

	return service
}

type AddRouteOptions struct {
	Path    *string
	Methods *[]awscdkapigatewayv2alpha.HttpMethod
}

func (s *ServiceApi) AddLambdaRoute(id string, handler awslambda.IFunction, options AddRouteOptions) *[]awscdkapigatewayv2alpha.HttpRoute {
	return s.HttpApi.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        options.Path,
		Methods:     options.Methods,
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(&id, handler, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{}),
	})
}

func (s *ServiceApi) AddServiceRoute(id string, service awsservicediscovery.IService, options AddRouteOptions) *[]awscdkapigatewayv2alpha.HttpRoute {
	return s.HttpApi.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:    options.Path,
		Methods: options.Methods,
		Integration: awscdkapigatewayv2integrationsalpha.NewHttpServiceDiscoveryIntegration(&id, service, &awscdkapigatewayv2integrationsalpha.HttpServiceDiscoveryIntegrationProps{
			VpcLink: *s.vpcLink,
		}),
	})
}

func (a ServiceApi) Connections() awsec2.Connections {
	return awsec2.NewConnections(&awsec2.ConnectionsProps{
		SecurityGroups: &[]awsec2.ISecurityGroup{},
	})
}
