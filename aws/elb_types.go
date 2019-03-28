package aws

import (
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// LoadBalancers - represents all load balancers
type LoadBalancers struct {
	Names []string
}

// GetAllResources - Gets all the elb names as a AwsResource
func (balancer LoadBalancers) GetAllResources(session *session.Session, region string, excludeAfter time.Time) (AwsResources, error) {
	elbNames, err := getAllElbInstances(session, region, excludeAfter)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	loadBalancers := LoadBalancers{
		Names: awsgo.StringValueSlice(elbNames),
	}

	return loadBalancers, nil
}

// ResourceName - the simple name of the aws resource
func (balancer LoadBalancers) ResourceName() string {
	return "elb"
}

// ResourceIdentifiers - The names of the load balancers
func (balancer LoadBalancers) ResourceIdentifiers() []string {
	return balancer.Names
}

func (balancer LoadBalancers) MaxBatchSize() int {
	// Tentative batch size to ensure AWS doesn't throttle
	return 200
}

// Nuke - nuke 'em all!!!
func (balancer LoadBalancers) Nuke(session *session.Session, identifiers []string) error {
	if err := nukeAllElbInstances(session, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}

type ElbDeleteError struct{}

func (e ElbDeleteError) Error() string {
	return "ELB was not deleted"
}
