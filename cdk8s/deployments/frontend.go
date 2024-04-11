package cdk8skit

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
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

	annotations := make(map[string]*string)

	if props.ClusterIssuer != nil {
		annotations["cert-manager.io/cluster-issuer"] = props.ClusterIssuer
	}

	ingress := cdk8splus28.NewIngress(scope, jsii.String("ingress"), &cdk8splus28.IngressProps{
		Metadata: &cdk8s.ApiObjectMetadata{
			Labels: &map[string]*string{
				"io.service": jsii.String(id),
			},
			Annotations: &annotations,
		},
		Rules: &[]*cdk8splus28.IngressRule{
			{
				Host:     host,
				Path:     jsii.String("/"),
				Backend:  cdk8splus28.IngressBackend_FromService(backend.Service, nil),
				PathType: cdk8splus28.HttpIngressPathType_PREFIX,
			},
		},
		Tls: &[]*cdk8splus28.IngressTls{
			{
				Hosts: &[]*string{
					host,
				},
				Secret: cdk8splus28.Secret_FromSecretName(scope, &id, jsii.String(fmt.Sprintf("%s-tls", id))),
			},
		},
	})

	return FrontendResource{
		Deployment: backend.Deployment,
		Service:    backend.Service,
		Ingress:    ingress,
	}
}
