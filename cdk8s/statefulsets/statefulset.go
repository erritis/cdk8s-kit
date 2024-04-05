package cdk8skit

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2"
)

type StatefulSetPort struct {
	Port          *float64
	ContainerPort *float64
}

type StatefulSetProps struct {
	Ports     *StatefulSetPort
	Network   *string
	Variables *map[*string]*string
	Claims    *[]*cdk8splus28.PersistentVolumeClaim
	Volumes   *map[*string]*cdk8splus28.Volume
	Liveness  cdk8splus28.Probe
}

func (props *StatefulSetProps) defaultProps() {

	if props.Ports == nil {
		props.Ports = &StatefulSetPort{}
	}
	if props.Ports.Port == nil {
		props.Ports.Port = jsii.Number(80)
	}
	if props.Ports.ContainerPort == nil {
		props.Ports.ContainerPort = jsii.Number(8080)
	}
}

func NewStatefulSet(
	scope constructs.Construct,
	id string,
	image string,
	props *StatefulSetProps,
) cdk8splus28.StatefulSet {

	props.defaultProps()

	container := cdk8splus28.NewContainer(&cdk8splus28.ContainerProps{
		Name:       jsii.String(id),
		Image:      jsii.String(image),
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
		Liveness: props.Liveness,
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

	statefulset := cdk8splus28.NewStatefulSet(
		scope,
		jsii.String("statefulset"),
		&cdk8splus28.StatefulSetProps{
			Replicas: jsii.Number(1),
			Service: cdk8splus28.NewService(
				scope,
				jsii.String("service"),
				&cdk8splus28.ServiceProps{
					Type: cdk8splus28.ServiceType_CLUSTER_IP,
					Ports: &[]*cdk8splus28.ServicePort{
						{
							Port:       props.Ports.Port,
							TargetPort: props.Ports.ContainerPort,
						},
					},
				},
			),
			SecurityContext: &cdk8splus28.PodSecurityContextProps{
				EnsureNonRoot: jsii.Bool(false),
			},
			PodMetadata: &cdk8s.ApiObjectMetadata{
				Labels: &labels,
			},
		},
	)

	statefulset.AttachContainer(container)

	if props.Claims != nil {
		claims := *props.Claims

		var jsonClaims []interface{}

		for _, claim := range claims {

			jsonClaims = append(jsonClaims, (*claim).ApiObject().ToJson())
		}

		statefulset.ApiObject().AddJsonPatch(
			cdk8s.JsonPatch_Add(
				jsii.String("/spec/volumeClaimTemplates"),
				&jsonClaims,
			),
		)
		for _, claim := range claims {
			claim_id := (*claim).Node().Id()
			scope.Node().TryRemoveChild(claim_id)
		}
	}

	statefulset.Metadata().AddLabel(jsii.String("io.service"), jsii.String(id))
	statefulset.Service().Metadata().AddLabel(jsii.String("io.service"), jsii.String(id))

	return statefulset
}
