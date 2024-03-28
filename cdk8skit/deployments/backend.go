package cdk8skit

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
	configs "github.com/erritis/cdk8skit/v2/cdk8skit/configs"
)

type TupleBackend struct {
	Deployment cdk8splus26.Deployment
	Service    cdk8splus26.Service
}

type BackendProps struct {
	PortConfig *configs.ServicePortConfig
	Network    *string
	Variables  *map[*string]*string
}

func (props *BackendProps) defaultProps() {

	if props.PortConfig == nil {
		props.PortConfig = &configs.ServicePortConfig{}
	}
	if props.PortConfig.Port == nil {
		props.PortConfig.Port = jsii.Number(80)
	}
	if props.PortConfig.ContainerPort == nil {
		props.PortConfig.ContainerPort = jsii.Number(8080)
	}
}

func NewBackend(
	scope constructs.Construct,
	id string,
	image *string,
	props *BackendProps,
) TupleBackend {

	props.defaultProps()

	container := cdk8splus26.NewContainer(&cdk8splus26.ContainerProps{
		Name:       jsii.String(id),
		Image:      image,
		PortNumber: props.PortConfig.ContainerPort,
		Resources: &cdk8splus26.ContainerResources{
			Cpu:              nil,
			EphemeralStorage: nil,
			Memory:           nil,
		},
		SecurityContext: &cdk8splus26.ContainerSecurityContextProps{
			ReadOnlyRootFilesystem: jsii.Bool(false),
			EnsureNonRoot:          jsii.Bool(false),
		},
	})

	for k, v := range *props.Variables {
		container.Env().AddVariable(k, cdk8splus26.EnvValue_FromValue(v))
	}

	labels := make(map[string]*string)

	if props.Network != nil {
		labels[*props.Network] = jsii.String("true")
	}

	deployment := cdk8splus26.NewDeployment(
		scope,
		jsii.String("deployment"),
		&cdk8splus26.DeploymentProps{
			Replicas: jsii.Number(1),
			SecurityContext: &cdk8splus26.PodSecurityContextProps{
				EnsureNonRoot: jsii.Bool(false),
			},
			PodMetadata: &cdk8s.ApiObjectMetadata{
				Labels: &labels,
			},
		},
	)

	deployment.AttachContainer(container)

	deployment.Metadata().AddLabel(jsii.String("io.service"), jsii.String(id))

	service := deployment.ExposeViaService(&cdk8splus26.DeploymentExposeViaServiceOptions{
		Name:        jsii.String(fmt.Sprintf("%s-service", id)),
		ServiceType: cdk8splus26.ServiceType_CLUSTER_IP,
		Ports: &[]*cdk8splus26.ServicePort{
			{
				Port:       props.PortConfig.Port,
				TargetPort: props.PortConfig.ContainerPort,
			},
		},
	})

	return TupleBackend{
		Deployment: deployment,
		Service:    service,
	}
}
