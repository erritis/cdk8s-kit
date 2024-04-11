package cdk8skit

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2"
)

type ClaimProps struct {
	StorageClassName *string
	Capacity         *cdk8s.Size
}

func (props *ClaimProps) defaultProps() {
	if props.StorageClassName == nil {
		props.StorageClassName = jsii.String("standard")
	}
	if props.Capacity == nil {
		capacity := cdk8s.Size_Gibibytes(jsii.Number(0.1))
		props.Capacity = &capacity
	}
}

func newClaim(scope constructs.Construct, id string, props *ClaimProps) cdk8splus28.PersistentVolumeClaim {

	props.defaultProps()

	claim := cdk8splus28.NewPersistentVolumeClaim(
		scope,
		jsii.String(id),
		&cdk8splus28.PersistentVolumeClaimProps{
			AccessModes: &[]cdk8splus28.PersistentVolumeAccessMode{
				cdk8splus28.PersistentVolumeAccessMode_READ_WRITE_ONCE,
			},
			Storage:          *props.Capacity,
			StorageClassName: props.StorageClassName,
		},
	)

	return claim
}
