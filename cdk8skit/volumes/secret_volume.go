package cdk8skit

import (
	"fmt"
	"os"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
)

type SecretVolumeProps struct {
	Encrypt *bool
}

func (props *SecretVolumeProps) defaultProps() {
	if props.Encrypt == nil {
		props.Encrypt = jsii.Bool(true)
	}
}

func NewSecretVolume(scope constructs.Construct, id string, name *string, value *string, props *SecretVolumeProps) cdk8splus26.Volume {

	props.defaultProps()

	secret := cdk8splus26.NewSecret(
		scope,
		jsii.String(id),
		&cdk8splus26.SecretProps{Type: jsii.String("Opaque")},
	)

	secret.AddStringData(name, value)

	if !*props.Encrypt {
		secret.ApiObject().AddJsonPatch(
			cdk8s.JsonPatch_Move(
				jsii.String("/stringData"),
				jsii.String("/data"),
			),
		)
	}

	volume := cdk8splus26.Volume_FromSecret(
		scope,
		name,
		secret,
		&cdk8splus26.SecretVolumeOptions{
			Name: name,
			Items: &map[string]*cdk8splus26.PathMapping{
				*name: {Path: name},
			},
		},
	)
	return volume
}

func readFileAsString(filename string) (string, error) {
	// Читаем содержимое файла в виде байтов
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// Преобразуем байты в строку
	content := string(bytes)
	return content, nil
}

func SecretVolume_FromFile(scope constructs.Construct, id string, path *string, filename *string) cdk8splus26.Volume {

	content, err := readFileAsString(*path)
	if err != nil {
		panic(fmt.Sprintf("Ошибка при чтении файла: %s", err))
	}
	return NewSecretVolume(scope, id, filename, &content, &SecretVolumeProps{
		Encrypt: jsii.Bool(false),
	})
}
