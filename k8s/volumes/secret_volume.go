package cdk8skit

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2/k8s"
)

type KubeSecretVolumeResource struct {
	Volume k8s.Volume
	Secret k8s.KubeSecret
}

type KubeSecretVolumeProps struct {
	Encrypt *bool
}

func (props *KubeSecretVolumeProps) defaultProps() {
	if props.Encrypt == nil {
		props.Encrypt = jsii.Bool(true)
	}
}

func NewKubeSecretVolume(scope constructs.Construct, id string, name *string, value *string, props *KubeSecretVolumeProps) KubeSecretVolumeResource {

	props.defaultProps()

	var data *map[string]*string

	var stringData *map[string]*string

	if *props.Encrypt {
		stringData = &map[string]*string{
			*name: value,
		}
		data = nil
	} else {
		stringData = nil
		data = &map[string]*string{
			*name: value,
		}
	}

	secret := k8s.NewKubeSecret(scope, &id, &k8s.KubeSecretProps{
		Type:       jsii.String("Opaque"),
		Data:       data,
		StringData: stringData,
		Immutable:  jsii.Bool(false),
	})

	volume := k8s.Volume{
		Name: name,
		Secret: &k8s.SecretVolumeSource{
			SecretName: jsii.String(fmt.Sprintf("%s-%s", *scope.Node().Id(), id)),
			Items: &[]*k8s.KeyToPath{
				{
					Key:  name,
					Path: name,
				},
			},
		},
	}

	return KubeSecretVolumeResource{
		Volume: volume,
		Secret: secret,
	}
}
