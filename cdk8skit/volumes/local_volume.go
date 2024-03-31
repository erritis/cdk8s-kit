package cdk8skit

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type LocalVolumeProps struct {
	StorageClassName *string
	Capacity         *cdk8s.Size
	Nodes            *[]string
}

func (props *LocalVolumeProps) defaultProps() {
	if props.Nodes == nil {
		props.Nodes = &[]string{"master-node"}
	}
}

func NewLocalVolume(
	scope constructs.Construct,
	id string,
	folder *string,
	props *LocalVolumeProps,
) VolumeResource {

	props.defaultProps()

	pvr := NewPersistentVolume(scope, id, &PersistentVolumeProps{
		StorageClassName: props.StorageClassName,
		Capacity:         props.Capacity,
	})

	pvr.PersistentVolume.ApiObject().AddJsonPatch(
		cdk8s.JsonPatch_Add(
			jsii.String("/spec/local"),
			&map[string]string{"path": *folder},
		),
	)

	pvr.PersistentVolume.ApiObject().AddJsonPatch(
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

	return VolumeResource{
		Volume: pvr.Volume,
		Claim:  pvr.Claim,
	}
}
