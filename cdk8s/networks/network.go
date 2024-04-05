package cdk8skit

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus28/v2"
)

func NewNetworkPolicy(
	scope constructs.Construct,
	id string,
	network string,
) cdk8splus28.NetworkPolicy {

	selector := cdk8splus28.Pods_Select(
		scope,
		jsii.String("selector"),
		&cdk8splus28.PodsSelectOptions{
			Labels: &map[string]*string{
				network: jsii.String("true"),
			},
		},
	)

	networkPolicy := cdk8splus28.NewNetworkPolicy(
		scope,
		jsii.String(id),
		&cdk8splus28.NetworkPolicyProps{
			Ingress: &cdk8splus28.NetworkPolicyTraffic{
				Rules: &[]*cdk8splus28.NetworkPolicyRule{
					{
						Peer: selector,
					},
				},
			},
			Selector: selector,
		},
	)

	return networkPolicy
}
