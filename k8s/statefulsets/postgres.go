package cdk8skit

import (
	"fmt"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2/k8s"
	volumes "github.com/erritis/cdk8skit/v4/k8s/volumes"
)

type KubePostgresResource struct {
	StatefulSet k8s.KubeStatefulSet
	Service     k8s.KubeService
	Props       KubePostgresProps
}

type KubePostgresVolumeSettings struct {
	PrefixSecretName     *string
	PrefixPersistentName *string
	StorageClassName     *string
	Capacity             *k8s.Quantity
	VolumeClaimTemplates *k8s.KubePersistentVolumeClaimProps
}

type KubePostgresDatabase struct {
	Name     *string
	Username *string
	Password *string
}

type KubePostgresPort struct {
	Port          *float64
	ContainerPort *float64
}

type KubePostgresProps struct {
	Image          *string
	Database       *KubePostgresDatabase
	Ports          *KubePostgresPort
	VolumeSettings *KubePostgresVolumeSettings
	Network        *string
	Liveness       *k8s.Probe
}

func (props *KubePostgresProps) defaultProps(id string) {
	if props.Image == nil {
		props.Image = jsii.String("postgres:latest")
	}

	props.defaultDbProps()

	props.defaultPortProps()

	props.defaultVolumeProps(id)

	props.defaultLivenessProps()
}

func (props *KubePostgresProps) defaultDbProps() {
	if props.Database == nil {
		props.Database = &KubePostgresDatabase{}
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

func (props *KubePostgresProps) defaultPortProps() {
	if props.Ports == nil {
		props.Ports = &KubePostgresPort{}
	}
	if props.Ports.Port == nil {
		props.Ports.Port = jsii.Number(5432)
	}
	if props.Ports.ContainerPort == nil {
		props.Ports.ContainerPort = jsii.Number(5432)
	}
}

func (props *KubePostgresProps) defaultVolumeProps(id string) {
	if props.VolumeSettings == nil {
		props.VolumeSettings = &KubePostgresVolumeSettings{}
	}
	if props.VolumeSettings.PrefixSecretName == nil {
		props.VolumeSettings.PrefixSecretName = jsii.String("postgres")
	}
	if props.VolumeSettings.PrefixPersistentName == nil {
		props.VolumeSettings.PrefixPersistentName = jsii.String(fmt.Sprintf("%s-claim", *props.VolumeSettings.PrefixSecretName))
	}
	if props.VolumeSettings.Capacity == nil {
		quantity := k8s.Quantity_FromString(jsii.String("0.1Gi"))
		props.VolumeSettings.Capacity = &quantity
	}
	if props.VolumeSettings.VolumeClaimTemplates == nil {
		props.VolumeSettings.VolumeClaimTemplates = &k8s.KubePersistentVolumeClaimProps{
			Metadata: &k8s.ObjectMeta{
				Name: props.VolumeSettings.PrefixPersistentName,
				Labels: &map[string]*string{
					"io.service": jsii.String(id),
				},
			},
			Spec: &k8s.PersistentVolumeClaimSpec{
				AccessModes: &[]*string{
					jsii.String("ReadWriteOnce"),
					jsii.String("ReadOnlyMany"),
				},
				Resources: &k8s.ResourceRequirements{
					Requests: &map[string]k8s.Quantity{
						"storage": *props.VolumeSettings.Capacity,
					},
				},
				StorageClassName: props.VolumeSettings.StorageClassName,
			},
		}
	}
}

func (props *KubePostgresProps) defaultLivenessProps() {
	if props.Liveness == nil {
		props.Liveness = &k8s.Probe{
			Exec: &k8s.ExecAction{
				Command: &[]*string{
					jsii.String("/bin/sh"),
					jsii.String("-c"),
					jsii.String("exec pg_isready -h 127.0.0.1"),
				},
			},
			FailureThreshold: jsii.Number(5),
			PeriodSeconds:    jsii.Number(5),
			TimeoutSeconds:   jsii.Number(5),
		}
	}
}

func NewKubePostgres(
	scope constructs.Construct,
	id string,
	props *KubePostgresProps,
) KubePostgresResource {

	props.defaultProps(id)

	db := volumes.NewKubeSecretVolume(
		scope, "name-secret",
		props.VolumeSettings.PrefixSecretName,
		props.Database.Name,
		&volumes.KubeSecretVolumeProps{},
	)

	dbUser := volumes.NewKubeSecretVolume(
		scope, "user-secret",
		jsii.String(fmt.Sprintf("%s-user", *props.VolumeSettings.PrefixSecretName)),
		props.Database.Username,
		&volumes.KubeSecretVolumeProps{},
	)

	dbPasswd := volumes.NewKubeSecretVolume(
		scope, "passwd-secret",
		jsii.String(fmt.Sprintf("%s-passwd", *props.VolumeSettings.PrefixSecretName)),
		props.Database.Password,
		&volumes.KubeSecretVolumeProps{},
	)

	statefulSetResource := NewKubeStatefulSet(
		scope,
		id,
		*props.Image,
		&KubeStatefulSetProps{
			Ports: &KubeStatefulSetPort{
				Port:          props.Ports.Port,
				ContainerPort: props.Ports.ContainerPort,
			},
			Network: props.Network,
			Variables: &map[string]*string{
				"POSTGRES_DB_FILE":       jsii.String(fmt.Sprintf("/run/secrets/%[1]s/%[1]s", *props.VolumeSettings.PrefixSecretName)),
				"POSTGRES_USER_FILE":     jsii.String(fmt.Sprintf("/run/secrets/%[1]s-user/%[1]s-user", *props.VolumeSettings.PrefixSecretName)),
				"POSTGRES_PASSWORD_FILE": jsii.String(fmt.Sprintf("/run/secrets/%[1]s-passwd/%[1]s-passwd", *props.VolumeSettings.PrefixSecretName)),
			},
			VolumeClaimTemplates: &map[string]*k8s.KubePersistentVolumeClaimProps{
				"/var/lib/postgresql/data": props.VolumeSettings.VolumeClaimTemplates,
			},
			Volumes: &map[string]*k8s.Volume{
				fmt.Sprintf("/run/secrets/%s", *props.VolumeSettings.PrefixSecretName):        &db.Volume,
				fmt.Sprintf("/run/secrets/%s-user", *props.VolumeSettings.PrefixSecretName):   &dbUser.Volume,
				fmt.Sprintf("/run/secrets/%s-passwd", *props.VolumeSettings.PrefixSecretName): &dbPasswd.Volume,
			},
			Liveness: props.Liveness,
		},
	)

	return KubePostgresResource{
		StatefulSet: statefulSetResource.StatefulSet,
		Service:     statefulSetResource.Service,
		Props:       *props,
	}
}
