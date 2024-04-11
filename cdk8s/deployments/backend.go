package cdk8skit

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2"
)

type BackendResource struct {
	Deployment cdk8splus28.Deployment
	Service    cdk8splus28.Service
}

type BackendPort struct {
	Port          *float64
	ContainerPort *float64
}

type BackendProps struct {
	Ports     *BackendPort
	Network   *string
	Variables *map[*string]*string
	Volumes   *map[*string]*cdk8splus28.Volume
}

func (props *BackendProps) defaultProps() {

	if props.Ports == nil {
		props.Ports = &BackendPort{}
	}
	if props.Ports.Port == nil {
		props.Ports.Port = jsii.Number(80)
	}
	if props.Ports.ContainerPort == nil {
		props.Ports.ContainerPort = jsii.Number(8080)
	}
	if props.Volumes == nil {
		props.Volumes = &map[*string]*cdk8splus28.Volume{}
	}
}

func NewBackend(
	scope constructs.Construct,
	id string,
	image *string,
	props *BackendProps,
) BackendResource {

	props.defaultProps()

	container := cdk8splus28.NewContainer(&cdk8splus28.ContainerProps{
		Name:       jsii.String(id),
		Image:      image,
		PortNumber: props.Ports.ContainerPort,
		Resources: &cdk8splus28.ContainerResources{
			Cpu:              nil,
			EphemeralStorage: nil,
			Memory:           nil,
		},
		SecurityContext: &cdk8splus28.ContainerSecurityContextProps{
			ReadOnlyRootFilesystem: jsii.Bool(false),
			EnsureNonRoot:          jsii.Bool(false),
		},
	})

	for k, v := range *props.Variables {
		container.Env().AddVariable(k, cdk8splus28.EnvValue_FromValue(v))
	}

	for path, volume := range *props.Volumes {
		var storage cdk8splus28.IStorage = *volume
		container.Mount(path, storage, nil)
	}

	labels := make(map[string]*string)

	if props.Network != nil {
		labels[*props.Network] = jsii.String("true")
	}

	deployment := cdk8splus28.NewDeployment(
		scope,
		jsii.String("deployment"),
		&cdk8splus28.DeploymentProps{
			Replicas: jsii.Number(1),
			SecurityContext: &cdk8splus28.PodSecurityContextProps{
				EnsureNonRoot: jsii.Bool(false),
			},
			PodMetadata: &cdk8s.ApiObjectMetadata{
				Labels: &labels,
			},
		},
	)

	deployment.AttachContainer(container)

	deployment.Metadata().AddLabel(jsii.String("io.service"), jsii.String(id))

	service := deployment.ExposeViaService(&cdk8splus28.DeploymentExposeViaServiceOptions{
		Name:        jsii.String(fmt.Sprintf("%s-service", id)),
		ServiceType: cdk8splus28.ServiceType_CLUSTER_IP,
		Ports: &[]*cdk8splus28.ServicePort{
			{
				Name:       jsii.String(fmt.Sprintf("%d", int(*props.Ports.Port))),
				Port:       props.Ports.Port,
				TargetPort: props.Ports.ContainerPort,
			},
		},
	})

	for _, volume := range *props.Volumes {
		deployment.AddVolume(*volume)
	}

	return BackendResource{
		Deployment: deployment,
		Service:    service,
	}
}
