package aws

import (
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// LoadBalancersV2 - represents all load balancers
type LoadBalancersV2 struct {
	Arns []string
}

// GetAllResources - Gets all the elbv2 names as a AwsResource
func (balancer LoadBalancersV2) GetAllResources(session *session.Session, region string, excludeAfter time.Time) (AwsResources, error) {
	elbv2Arns, err := getAllElbv2Instances(session, region, excludeAfter)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	loadBalancersV2 := LoadBalancersV2{
		Arns: awsgo.StringValueSlice(elbv2Arns),
	}

	return loadBalancersV2, nil
}

// ResourceName - the simple name of the aws resource
func (balancer LoadBalancersV2) ResourceName() string {
	return "elbv2"
}

func (balancer LoadBalancersV2) MaxBatchSize() int {
	// Tentative batch size to ensure AWS doesn't throttle
	return 200
}

// ResourceIdentifiers - The arns of the load balancers
func (balancer LoadBalancersV2) ResourceIdentifiers() []string {
	return balancer.Arns
}

// Nuke - nuke 'em all!!!
func (balancer LoadBalancersV2) Nuke(session *session.Session, identifiers []string) error {
	if err := nukeAllElbv2Instances(session, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
