package cdk8skit

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2/k8s"
)

type KubeLocalVolumeProps struct {
	StorageClassName *string
	Capacity         *k8s.Quantity
	Nodes            *[]string
}

func (props *KubeLocalVolumeProps) defaultProps() {
	if props.Nodes == nil {
		props.Nodes = &[]string{"master-node"}
	}
}

func NewKubeLocalVolume(scope constructs.Construct, id string, claim_name *string, folder *string, props *KubeLocalVolumeProps) k8s.KubePersistentVolume {

	props.defaultProps()

	volume := NewKubePersistentVolume(scope, id, claim_name, &KubePersistentVolumeProps{
		StorageClassName: props.StorageClassName,
		Capacity:         props.Capacity,
	})

	volume.AddJsonPatch(
		cdk8s.JsonPatch_Add(
			jsii.String("/spec/local"),
			&map[string]string{"path": *folder},
		),
	)

	volume.AddJsonPatch(
		cdk8s.JsonPatch_Add(
			jsii.String("/spec/nodeAffinity"),
			&map[string]interface{}{
				"required": &map[string]interface{}{
					"nodeSelectorTerms": &[]interface{}{
						&map[string]interface{}{
							"matchExpressions": &[]interface{}{
								&map[string]interface{}{
									"key":      jsii.String("kubernetes.io/hostname"),
									"operator": jsii.String("In"),
									"values":   props.Nodes,
								},
							},
						},
					},
				},
			},
		),
	)

	return volume
}
