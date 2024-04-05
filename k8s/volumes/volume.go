package cdk8skit

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2/k8s"
)

type KubeVolumeResource struct {
	Volume k8s.KubePersistentVolume
	Claim  k8s.KubePersistentVolumeClaim
}

func NewKubeVolume(scope constructs.Construct, id string) KubeVolumeResource {

	claim_id := fmt.Sprintf("%s-claim", id)

	claim_name := fmt.Sprintf("%s-%s", *scope.Node().Id(), claim_id)

	claim := k8s.NewKubePersistentVolumeClaim(
		scope,
		&claim_id,
		&k8s.KubePersistentVolumeClaimProps{
			Metadata: &k8s.ObjectMeta{
				Name: &claim_name,
			},
		},
	)

	volume := k8s.NewKubePersistentVolume(
		scope,
		&id,
		&k8s.KubePersistentVolumeProps{
			Spec: &k8s.PersistentVolumeSpec{
				ClaimRef: &k8s.ObjectReference{
					Name: &claim_id,
				},
			},
			Metadata: &k8s.ObjectMeta{
				Name: &claim_name,
			},
		},
	)

	return KubeVolumeResource{
		Volume: volume,
		Claim:  claim,
	}
}
