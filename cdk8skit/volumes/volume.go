package cdk8skit

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
)

type VolumeResource struct {
	Volume cdk8splus26.Volume
	Claim  cdk8splus26.IPersistentVolumeClaim
}

type VolumeProps struct {
	StorageClassName *string
	Capacity         *cdk8s.Size
}

func NewVolume(scope constructs.Construct, id string, props *VolumeProps) VolumeResource {

	claim_id := fmt.Sprintf("%s-claim", id)

	claim_name := fmt.Sprintf("%s-%s", *scope.Node().Id(), claim_id)

	claim := newClaim(scope, claim_id, &ClaimProps{
		StorageClassName: props.StorageClassName,
		Capacity:         props.Capacity,
	})

	volume := cdk8splus26.Volume_FromPersistentVolumeClaim(
		scope,
		jsii.String(id),
		claim,
		&cdk8splus26.PersistentVolumeClaimVolumeOptions{
			Name: jsii.String(claim_name),
		},
	)

	return VolumeResource{
		Volume: volume,
		Claim:  claim,
	}
}
