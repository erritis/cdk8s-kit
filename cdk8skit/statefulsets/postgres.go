package cdk8skit

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
	configs "github.com/erritis/cdk8skit/v2/cdk8skit/configs"
	volumes "github.com/erritis/cdk8skit/v2/cdk8skit/volumes"
)

type VolumeDefaultConfig struct {
	PrefixPersistentName *string
	VolumeProps          *volumes.VolumeProps
}

type DBConfig struct {
	Name     *string
	Username *string
	Password *string
}

type PostgresProps struct {
	Image               *string
	PrefixSecretName    *string
	DBConfig            *DBConfig
	PortConfig          *configs.ServicePortConfig
	VolumeDefaultConfig *VolumeDefaultConfig
	Volume              *cdk8splus26.Volume
	Network             *string
}

func (props *PostgresProps) defaultProps(scope constructs.Construct) {
	if props.Image == nil {
		props.Image = jsii.String("postgres:latest")
	}
	if props.PrefixSecretName == nil {
		props.PrefixSecretName = jsii.String("postgres")
	}
	if props.DBConfig == nil {
		props.DBConfig = &DBConfig{}
	}
	if props.DBConfig.Name == nil {
		props.DBConfig.Name = jsii.String("postgres")
	}
	if props.DBConfig.Username == nil {
		props.DBConfig.Username = jsii.String("postgres")
	}
	if props.DBConfig.Password == nil {
		props.DBConfig.Password = jsii.String("postgres")
	}
	if props.PortConfig == nil {
		props.PortConfig = &configs.ServicePortConfig{}
	}
	if props.PortConfig.Port == nil {
		props.PortConfig.Port = jsii.Number(5432)
	}
	if props.PortConfig.ContainerPort == nil {
		props.PortConfig.ContainerPort = jsii.Number(5432)
	}
	if props.VolumeDefaultConfig == nil {
		props.VolumeDefaultConfig = &VolumeDefaultConfig{}
	}
	if props.VolumeDefaultConfig.PrefixPersistentName == nil {
		props.VolumeDefaultConfig.PrefixPersistentName = jsii.String("persistent-volume")
	}
	if props.Volume == nil {
		disk := volumes.NewVolume(
			scope,
			*props.VolumeDefaultConfig.PrefixPersistentName,
			props.VolumeDefaultConfig.VolumeProps,
		)
		props.Volume = &disk.Volume
	}
}

func NewPostgres(
	scope constructs.Construct,
	id string,
	props *PostgresProps,
) cdk8splus26.StatefulSet {

	props.defaultProps(scope)

	db := volumes.NewSecretVolume(
		scope, "db-secret",
		props.PrefixSecretName,
		props.DBConfig.Name,
	)

	dbUser := volumes.NewSecretVolume(
		scope, "user-secret",
		jsii.String(fmt.Sprintf("%s-user", *props.PrefixSecretName)),
		props.DBConfig.Username,
	)

	dbPasswd := volumes.NewSecretVolume(
		scope, "passwd-secret",
		jsii.String(fmt.Sprintf("%s-passwd", *props.PrefixSecretName)),
		props.DBConfig.Password,
	)

	postgres := NewStatefulSet(
		scope,
		id,
		*props.Image,
		&StatefulSetProps{
			PortConfig: props.PortConfig,
			Network:    props.Network,
			Variables: &map[*string]*string{
				jsii.String("POSTGRES_DB_FILE"):       jsii.String(fmt.Sprintf("/run/secrets/%[1]s/%[1]s", *props.PrefixSecretName)),
				jsii.String("POSTGRES_USER_FILE"):     jsii.String(fmt.Sprintf("/run/secrets/%[1]s-user/%[1]s-user", *props.PrefixSecretName)),
				jsii.String("POSTGRES_PASSWORD_FILE"): jsii.String(fmt.Sprintf("/run/secrets/%[1]s-passwd/%[1]s-passwd", *props.PrefixSecretName)),
			},
			Volumes: &map[*string]*cdk8splus26.Volume{
				jsii.String("/var/lib/postgresql/data"):                                     props.Volume,
				jsii.String(fmt.Sprintf("/run/secrets/%s", *props.PrefixSecretName)):        &db,
				jsii.String(fmt.Sprintf("/run/secrets/%s-user", *props.PrefixSecretName)):   &dbUser,
				jsii.String(fmt.Sprintf("/run/secrets/%s-passwd", *props.PrefixSecretName)): &dbPasswd,
			},
		},
	)

	return postgres
}
