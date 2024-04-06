package cdk8skit

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2/k8s"
)

type KubePersistentVolumeProps struct {
	StorageClassName *string
	Capacity         *k8s.Quantity
}

func (props *KubePersistentVolumeProps) defaultProps() {
	if props.StorageClassName == nil {
		props.StorageClassName = jsii.String("standard")
	}
	if props.Capacity == nil {
		capacity := k8s.Quantity_FromString(jsii.String("0.1Gi"))
		props.Capacity = &capacity
	}
}

func NewKubePersistentVolume(scope constructs.Construct, id string, claim_name *string, props *KubePersistentVolumeProps) k8s.KubePersistentVolume {

	props.defaultProps()

	volume := k8s.NewKubePersistentVolume(
		scope,
		&id,
		&k8s.KubePersistentVolumeProps{
			Spec: &k8s.PersistentVolumeSpec{
				ClaimRef: &k8s.ObjectReference{
					Name: claim_name,
				},
				AccessModes: &[]*string{
					jsii.String("ReadWriteOnce"),
					jsii.String("ReadOnlyMany"),
				},
				Capacity: &map[string]k8s.Quantity{
					"storage": *props.Capacity,
				},
				StorageClassName: props.StorageClassName,
			},
		},
	)

	return volume
}
