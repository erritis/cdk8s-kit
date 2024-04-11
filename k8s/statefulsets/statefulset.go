package cdk8skit

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2/k8s"
)

type KubeStatefulSetResource struct {
	StatefulSet k8s.KubeStatefulSet
	Service     k8s.KubeService
}

type KubeStatefulSetPort struct {
	Port          *float64
	ContainerPort *float64
}

type KubeStatefulSetProps struct {
	Ports                *KubeStatefulSetPort
	Network              *string
	Variables            *map[string]*string
	VolumeClaimTemplates *map[string]*k8s.KubePersistentVolumeClaimProps
	Volumes              *map[string]*k8s.Volume
	Liveness             *k8s.Probe
}

func (props *KubeStatefulSetProps) defaultProps() {

	if props.Ports == nil {
		props.Ports = &KubeStatefulSetPort{}
	}
	if props.Ports.Port == nil {
		props.Ports.Port = jsii.Number(80)
	}
	if props.Ports.ContainerPort == nil {
		props.Ports.ContainerPort = jsii.Number(8080)
	}
}

func NewKubeStatefulSet(
	scope constructs.Construct,
	id string,
	image string,
	props *KubeStatefulSetProps,
) KubeStatefulSetResource {

	props.defaultProps()

	labels := make(map[string]*string)

	labels["io.service"] = jsii.String(id)

	if props.Network != nil {
		labels[*props.Network] = jsii.String("true")
	}

	service := k8s.NewKubeService(
		scope,
		jsii.String("service"),
		&k8s.KubeServiceProps{
			Metadata: &k8s.ObjectMeta{
				Labels: &labels,
			},
			Spec: &k8s.ServiceSpec{
				Selector: &map[string]*string{
					"io.service": labels["io.service"],
				},
				Ports: &[]*k8s.ServicePort{
					{
						Name:       jsii.String(fmt.Sprintf("%d", int(*props.Ports.Port))),
						Port:       props.Ports.Port,
						TargetPort: k8s.IntOrString_FromNumber(props.Ports.ContainerPort),
					},
				},
				Type: jsii.String("ClusterIP"),
			},
		},
	)

	variables := []*k8s.EnvVar{}

	for k, v := range *props.Variables {
		variables = append(variables, &k8s.EnvVar{
			Name:  &k,
			Value: v,
		})
	}

	mounts := []*k8s.VolumeMount{}

	for path, volume := range *props.Volumes {
		mounts = append(mounts, &k8s.VolumeMount{
			MountPath: &path,
			Name:      volume.Name,
		})
	}

	var volumeClaimTemplates []*k8s.KubePersistentVolumeClaimProps

	for path, claim := range *props.VolumeClaimTemplates {
		mounts = append(mounts, &k8s.VolumeMount{
			MountPath: &path,
			Name:      claim.Metadata.Name,
		})
		volumeClaimTemplates = append(volumeClaimTemplates, claim)
	}

	volumes := []*k8s.Volume{}

	for _, volume := range *props.Volumes {
		volumes = append(volumes, volume)
	}

	statefulset := k8s.NewKubeStatefulSet(
		scope,
		jsii.String("statefulset"),
		&k8s.KubeStatefulSetProps{
			Metadata: &k8s.ObjectMeta{
				Name: jsii.String(fmt.Sprintf("%s-statefulset", id)),
				Labels: &map[string]*string{
					"io.service": labels["io.service"],
				},
			},
			Spec: &k8s.StatefulSetSpec{
				Replicas: jsii.Number(1),
				Selector: &k8s.LabelSelector{
					MatchLabels: &map[string]*string{
						"io.service": labels["io.service"],
					},
				},
				ServiceName:          service.Name(),
				VolumeClaimTemplates: &volumeClaimTemplates,
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels: &labels,
					},
					Spec: &k8s.PodSpec{
						Containers: &[]*k8s.Container{
							{
								Name:      jsii.String(fmt.Sprintf("%s-statefulset-pod", id)),
								Image:     jsii.String(image),
								Resources: &k8s.ResourceRequirements{},
								Ports: &[]*k8s.ContainerPort{
									{
										ContainerPort: props.Ports.ContainerPort,
									},
								},
								LivenessProbe: props.Liveness,
								SecurityContext: &k8s.SecurityContext{
									RunAsNonRoot: jsii.Bool(false),
								},
								Env:          &variables,
								VolumeMounts: &mounts,
							},
						},
						SecurityContext: &k8s.PodSecurityContext{
							RunAsNonRoot: jsii.Bool(false),
						},
						Volumes: &volumes,
					},
				},
			},
		},
	)

	return KubeStatefulSetResource{
		StatefulSet: statefulset,
		Service:     service,
	}
}
