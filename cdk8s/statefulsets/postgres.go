package cdk8skit

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2"
	volumes "github.com/erritis/cdk8skit/v4/cdk8s/volumes"
)

type PostgresPort struct {
	Port          *float64
	ContainerPort *float64
}

type PostgresVolumeSettings struct {
	PrefixSecretName     *string
	PrefixPersistentName *string
	StorageClassName     *string
	Capacity             *cdk8s.Size
	Volume               *cdk8splus28.Volume
	Claim                *cdk8splus28.PersistentVolumeClaim
}

type PostgresDatabase struct {
	Name     *string
	Username *string
	Password *string
}

type PostgresProps struct {
	Image          *string
	Database       *PostgresDatabase
	VolumeSettings *PostgresVolumeSettings
	Ports          *PostgresPort
	Network        *string
	Liveness       cdk8splus28.Probe
}

func (props *PostgresProps) defaultProps(scope constructs.Construct) {
	if props.Image == nil {
		props.Image = jsii.String("postgres:latest")
	}
	if props.VolumeSettings.PrefixSecretName == nil {
		props.VolumeSettings.PrefixSecretName = jsii.String("postgres")
	}

	props.defaultDbProps()

	props.defaultPortProps()

	props.defaultVolumeProps(scope)

	props.defaultLivenessProps()
}

func (props *PostgresProps) defaultDbProps() {
	if props.Database == nil {
		props.Database = &PostgresDatabase{}
	}
	if props.Database.Name == nil {
		props.Database.Name = jsii.String("postgres")
	}
	if props.Database.Username == nil {
		props.Database.Username = jsii.String("postgres")
	}
	if props.Database.Password == nil {
		props.Database.Password = jsii.String("postgres")
	}
}

func (props *PostgresProps) defaultPortProps() {
	if props.Ports == nil {
		props.Ports = &PostgresPort{}
	}
	if props.Ports.Port == nil {
		props.Ports.Port = jsii.Number(5432)
	}
	if props.Ports.ContainerPort == nil {
		props.Ports.ContainerPort = jsii.Number(5432)
	}
}

func (props *PostgresProps) defaultVolumeProps(scope constructs.Construct) {
	if props.VolumeSettings == nil {
		props.VolumeSettings = &PostgresVolumeSettings{}
	}
	if props.VolumeSettings.PrefixPersistentName == nil {
		props.VolumeSettings.PrefixPersistentName = jsii.String("persistent-volume")
	}
	if props.VolumeSettings.Volume == nil && props.VolumeSettings.Claim == nil {
		volumeResource := volumes.NewVolume(
			scope,
			*props.VolumeSettings.PrefixPersistentName,
			&volumes.VolumeProps{
				StorageClassName: props.VolumeSettings.StorageClassName,
				Capacity:         props.VolumeSettings.Capacity,
			},
		)
		props.VolumeSettings.Volume = &volumeResource.Volume
		props.VolumeSettings.Claim = &volumeResource.Claim
	}
}

func (props *PostgresProps) defaultLivenessProps() {
	if props.Liveness == nil {
		props.Liveness = cdk8splus28.Probe_FromCommand(
			&[]*string{
				jsii.String("/bin/sh"),
				jsii.String("-c"),
				jsii.String("exec pg_isready -h 127.0.0.1"),
			},
			&cdk8splus28.CommandProbeOptions{
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
) cdk8splus28.StatefulSet {

	props.defaultProps(scope)

	db := volumes.NewSecretVolume(
		scope, "name-secret",
		props.VolumeSettings.PrefixSecretName,
		props.Database.Name,
		&volumes.SecretVolumeProps{},
	)

	dbUser := volumes.NewSecretVolume(
		scope, "user-secret",
		jsii.String(fmt.Sprintf("%s-user", *props.VolumeSettings.PrefixSecretName)),
		props.Database.Username,
		&volumes.SecretVolumeProps{},
	)

	dbPasswd := volumes.NewSecretVolume(
		scope, "passwd-secret",
		jsii.String(fmt.Sprintf("%s-passwd", *props.VolumeSettings.PrefixSecretName)),
		props.Database.Password,
		&volumes.SecretVolumeProps{},
	)

	postgres := NewStatefulSet(
		scope,
		id,
		*props.Image,
		&StatefulSetProps{
			Ports: &StatefulSetPort{
				Port:          props.Ports.Port,
				ContainerPort: props.Ports.ContainerPort,
			},
			Network: props.Network,
			Variables: &map[*string]*string{
				jsii.String("POSTGRES_DB_FILE"):       jsii.String(fmt.Sprintf("/run/secrets/%[1]s/%[1]s", *props.VolumeSettings.PrefixSecretName)),
				jsii.String("POSTGRES_USER_FILE"):     jsii.String(fmt.Sprintf("/run/secrets/%[1]s-user/%[1]s-user", *props.VolumeSettings.PrefixSecretName)),
				jsii.String("POSTGRES_PASSWORD_FILE"): jsii.String(fmt.Sprintf("/run/secrets/%[1]s-passwd/%[1]s-passwd", *props.VolumeSettings.PrefixSecretName)),
			},
			Claims: &[]*cdk8splus28.PersistentVolumeClaim{
				props.VolumeSettings.Claim,
			},
			Volumes: &map[*string]*cdk8splus28.Volume{
				jsii.String("/var/lib/postgresql/data"):                                                    props.VolumeSettings.Volume,
				jsii.String(fmt.Sprintf("/run/secrets/%s", *props.VolumeSettings.PrefixSecretName)):        &db,
				jsii.String(fmt.Sprintf("/run/secrets/%s-user", *props.VolumeSettings.PrefixSecretName)):   &dbUser,
				jsii.String(fmt.Sprintf("/run/secrets/%s-passwd", *props.VolumeSettings.PrefixSecretName)): &dbPasswd,
			},
			Liveness: props.Liveness,
		},
	)

	return postgres
}
