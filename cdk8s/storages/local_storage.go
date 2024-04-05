package cdk8skit

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2/k8s"
)

type LocalStorageProps struct {
	IsDefault *bool
}

func (props *LocalStorageProps) defaultProps() {
	if props.IsDefault == nil {
		props.IsDefault = jsii.Bool(false)
	}
}

func NewLocalStorage(scope constructs.Construct, id string, props *LocalStorageProps) k8s.KubeStorageClass {

	props.defaultProps()

	annotations := make(map[string]*string)

	annotations["storageclass.kubernetes.io/is-default-class"] = jsii.String(fmt.Sprintf("%t", *props.IsDefault))

	metadata := k8s.ObjectMeta{
		Name:        jsii.String(id),
		Annotations: &annotations,
	}

	storageProps := k8s.KubeStorageClassProps{
		Provisioner:       jsii.String("kubernetes.io/no-provisioner"),
		Metadata:          &metadata,
		VolumeBindingMode: jsii.String("WaitForFirstConsumer"),
	}

	storage := k8s.NewKubeStorageClass(scope, jsii.String(id), &storageProps)

	storage.AddJsonPatch(cdk8s.JsonPatch_Replace(jsii.String("/metadata/namespace"), new(*string)))

	return storage
}
