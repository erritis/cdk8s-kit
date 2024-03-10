package cdk8skit

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
)

type TupleFrontend struct {
	Deployment cdk8splus26.Deployment
	Service    cdk8splus26.Service
	Ingress    cdk8splus26.Ingress
}

func NewFrontend(
	scope constructs.Construct,
	id string,
	host *string,
	cluster_issuer *string,
	image *string,
	port *float64,
	containerPort *float64,
	network string,
	variables *map[*string]*string,
) TupleFrontend {
	backend := NewBackend(scope, id, image, port, containerPort, network, variables)

	ingress := cdk8splus26.NewIngress(scope, jsii.String("ingress"), nil)

	ingress.Metadata().AddLabel(jsii.String("io.service"), jsii.String(id))
	ingress.Metadata().AddAnnotation(jsii.String("cert-manager.io/cluster-issuer"), cluster_issuer)

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
