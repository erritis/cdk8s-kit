package cdk8skit

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus26/v2"
	configs "github.com/erritis/cdk8skit/v3/cdk8skit/configs"
	volumes "github.com/erritis/cdk8skit/v3/cdk8skit/volumes"
)

type VolumeConfig struct {
	PrefixPersistentName *string
	StorageClassName     *string
	Capacity             *cdk8s.Size
	Volume               *cdk8splus26.Volume
}

type DBConfig struct {
	Name     *string
	Username *string
	Password *string
}

type PostgresProps struct {
	Image            *string
	PrefixSecretName *string
	DBConfig         *DBConfig
	PortConfig       *configs.ServicePortConfig
	VolumeConfig     *VolumeConfig
	Network          *string
	Liveness         cdk8splus26.Probe
}

func (props *PostgresProps) defaultProps(scope constructs.Construct) {
	if props.Image == nil {
		props.Image = jsii.String("postgres:latest")
	}
	if props.PrefixSecretName == nil {
		props.PrefixSecretName = jsii.String("postgres")
	}

	props.defaultDbProps()

	props.defaultPortProps()

	props.defaultVolumeProps(scope)

	props.defaultLivenessProps()
}

func (props *PostgresProps) defaultDbProps() {
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
}

func (props *PostgresProps) defaultPortProps() {
	if props.PortConfig == nil {
		props.PortConfig = &configs.ServicePortConfig{}
	}
	if props.PortConfig.Port == nil {
		props.PortConfig.Port = jsii.Number(5432)
	}
	if props.PortConfig.ContainerPort == nil {
		props.PortConfig.ContainerPort = jsii.Number(5432)
	}
}

func (props *PostgresProps) defaultVolumeProps(scope constructs.Construct) {
	if props.VolumeConfig == nil {
		props.VolumeConfig = &VolumeConfig{}
	}
	if props.VolumeConfig.PrefixPersistentName == nil {
		props.VolumeConfig.PrefixPersistentName = jsii.String("persistent-volume")
	}
	if props.VolumeConfig.Volume == nil {
		hv := volumes.NewVolume(
			scope,
			*props.VolumeConfig.PrefixPersistentName,
			&volumes.VolumeProps{
				StorageClassName: props.VolumeConfig.StorageClassName,
				Capacity:         props.VolumeConfig.Capacity,
			},
		)
		props.VolumeConfig.Volume = &hv.Volume
	}
}

func (props *PostgresProps) defaultLivenessProps() {
	if props.Liveness == nil {
		props.Liveness = cdk8splus26.Probe_FromCommand(
			&[]*string{
				jsii.String("/bin/sh"),
				jsii.String("-c"),
				jsii.String("exec pg_isready -h 127.0.0.1"),
			},
			&cdk8splus26.CommandProbeOptions{
				FailureThreshold: jsii.Number(5),
				PeriodSeconds:    cdk8s.Duration_Seconds(jsii.Number(5)),
				TimeoutSeconds:   cdk8s.Duration_Seconds(jsii.Number(5)),
			},
		)
	}
}

func NewPostgres(
	scope constructs.Construct,
	id string,
	props *PostgresProps,
) cdk8splus26.StatefulSet {

	props.defaultProps(scope)

	db := volumes.NewSecretVolume(
		scope, "name-secret",
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
				jsii.String("/var/lib/postgresql/data"):                                     props.VolumeConfig.Volume,
				jsii.String(fmt.Sprintf("/run/secrets/%s", *props.PrefixSecretName)):        &db,
				jsii.String(fmt.Sprintf("/run/secrets/%s-user", *props.PrefixSecretName)):   &dbUser,
				jsii.String(fmt.Sprintf("/run/secrets/%s-passwd", *props.PrefixSecretName)): &dbPasswd,
			},
			Liveness: props.Liveness,
		},
	)

	return postgres
}
