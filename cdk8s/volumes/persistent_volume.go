package cdk8skit

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2"
)

type PersistentVolumeResource struct {
	PersistentVolume cdk8splus28.PersistentVolume
	Volume           cdk8splus28.Volume
	Claim            cdk8splus28.PersistentVolumeClaim
}

type PersistentVolumeProps struct {
	StorageClassName *string
	Capacity         *cdk8s.Size
}

func (props *PersistentVolumeProps) defaultProps() {
	if props.StorageClassName == nil {
		props.StorageClassName = jsii.String("standard")
	}
	if props.Capacity == nil {
		capacity := cdk8s.Size_Gibibytes(jsii.Number(0.1))
		props.Capacity = &capacity
	}
}

func NewPersistentVolume(scope constructs.Construct, id string, props *PersistentVolumeProps) PersistentVolumeResource {

	props.defaultProps()

	claim_id := fmt.Sprintf("%s-claim", id)

	claim := newClaim(scope, claim_id, &ClaimProps{
		StorageClassName: props.StorageClassName,
		Capacity:         props.Capacity,
	})

	claim_name := fmt.Sprintf("%s-%s", *scope.Node().Id(), claim_id)

	volume := cdk8splus28.Volume_FromPersistentVolumeClaim(
		scope,
		jsii.String(fmt.Sprintf("%s-ref", id)),
		claim,
		&cdk8splus28.PersistentVolumeClaimVolumeOptions{
			Name: jsii.String(claim_name),
		},
	)

	persistentVolume := cdk8splus28.NewPersistentVolume(
		scope,
		jsii.String(id),
		&cdk8splus28.PersistentVolumeProps{
			VolumeMode: cdk8splus28.PersistentVolumeMode_FILE_SYSTEM,
			AccessModes: &[]cdk8splus28.PersistentVolumeAccessMode{
				cdk8splus28.PersistentVolumeAccessMode_READ_WRITE_ONCE,
				cdk8splus28.PersistentVolumeAccessMode_READ_ONLY_MANY,
			},
			ReclaimPolicy:    cdk8splus28.PersistentVolumeReclaimPolicy_RETAIN,
			Storage:          *props.Capacity,
			StorageClassName: props.StorageClassName,
			Claim:            claim,
		},
	)

	return PersistentVolumeResource{
		PersistentVolume: persistentVolume,
		Volume:           volume,
		Claim:            claim,
	}
}
