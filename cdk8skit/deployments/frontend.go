package cdk8skit

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
	configs "github.com/erritis/cdk8skit/v3/cdk8skit/configs"
)

type TupleFrontend struct {
	Deployment cdk8splus26.Deployment
	Service    cdk8splus26.Service
	Ingress    cdk8splus26.Ingress
}

type FrontendProps struct {
	PortConfig    *configs.ServicePortConfig
	Network       *string
	Variables     *map[*string]*string
	Volumes       *map[*string]*cdk8splus26.Volume
	ClusterIssuer *string
}

func (props *FrontendProps) defaultProps() {

	if props.PortConfig == nil {
		props.PortConfig = &configs.ServicePortConfig{}
	}
}

func NewFrontend(
	scope constructs.Construct,
	id string,
	host *string,
	image *string,
	props *FrontendProps,
) TupleFrontend {

	props.defaultProps()

	backend := NewBackend(scope, id, image, &BackendProps{
		PortConfig: props.PortConfig,
		Network:    props.Network,
		Variables:  props.Variables,
		Volumes:    props.Volumes,
	})

	ingress := cdk8splus26.NewIngress(scope, jsii.String("ingress"), nil)

	ingress.Metadata().AddLabel(jsii.String("io.service"), jsii.String(id))

	if props.ClusterIssuer != nil {
		ingress.Metadata().AddAnnotation(jsii.String("cert-manager.io/cluster-issuer"), props.ClusterIssuer)
	}

	ingress.AddHostRule(
		host,
		jsii.String("/"),
		cdk8splus26.IngressBackend_FromService(backend.Service, nil),
		cdk8splus26.HttpIngressPathType_PREFIX,
	)

	return TupleFrontend{
		Deployment: backend.Deployment,
		Service:    backend.Service,
		Ingress:    ingress,
	}
}
