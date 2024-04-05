package cdk8skit

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2"
)

type FrontendResource struct {
	Deployment cdk8splus28.Deployment
	Service    cdk8splus28.Service
	Ingress    cdk8splus28.Ingress
}

type FrontendPort struct {
	Port          *float64
	ContainerPort *float64
}

type FrontendProps struct {
	Ports         *FrontendPort
	Network       *string
	Variables     *map[*string]*string
	Volumes       *map[*string]*cdk8splus28.Volume
	ClusterIssuer *string
}

func (props *FrontendProps) defaultProps() {

	if props.Ports == nil {
		props.Ports = &FrontendPort{}
	}
}

func NewFrontend(
	scope constructs.Construct,
	id string,
	host *string,
	image *string,
	props *FrontendProps,
) FrontendResource {

	props.defaultProps()

	backend := NewBackend(scope, id, image, &BackendProps{
		Ports: &BackendPort{
			Port:          props.Ports.Port,
			ContainerPort: props.Ports.ContainerPort,
		},
		Network:   props.Network,
		Variables: props.Variables,
		Volumes:   props.Volumes,
	})

	ingress := cdk8splus28.NewIngress(scope, jsii.String("ingress"), nil)

	ingress.Metadata().AddLabel(jsii.String("io.service"), jsii.String(id))

	if props.ClusterIssuer != nil {
		ingress.Metadata().AddAnnotation(jsii.String("cert-manager.io/cluster-issuer"), props.ClusterIssuer)
	}

	ingress.AddHostRule(
		host,
		jsii.String("/"),
		cdk8splus28.IngressBackend_FromService(backend.Service, nil),
		cdk8splus28.HttpIngressPathType_PREFIX,
	)

	return FrontendResource{
		Deployment: backend.Deployment,
		Service:    backend.Service,
		Ingress:    ingress,
	}
}
