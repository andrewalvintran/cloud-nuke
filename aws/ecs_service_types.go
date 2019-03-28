package aws

import (
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// ECSServices - Represents all ECS services found in a region
type ECSServices struct {
	Services          []string
	ServiceClusterMap map[string]string
}

// GetAllResources - Gets all the ECS services as a AwsResource
func (services ECSServices) GetAllResources(session *session.Session, region string, excludeAfter time.Time) (AwsResources, error) {
	clusterArns, err := getAllEcsClusters(session)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}
	serviceArns, serviceClusterMap, err := getAllEcsServices(session, clusterArns, excludeAfter)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	ecsServices := ECSServices{
		Services:          awsgo.StringValueSlice(serviceArns),
		ServiceClusterMap: serviceClusterMap,
	}

	return ecsServices, nil
}

// ResourceName - The simple name of the aws resource
func (services ECSServices) ResourceName() string {
	return "ecsserv"
}

// ResourceIdentifiers - The ARNs of the collected ECS services
func (services ECSServices) ResourceIdentifiers() []string {
	return services.Services
}

func (services ECSServices) MaxBatchSize() int {
	return 200
}

// Nuke - nuke all ECS service resources
func (services ECSServices) Nuke(awsSession *session.Session, identifiers []string) error {
	if err := nukeAllEcsServices(awsSession, services.ServiceClusterMap, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}
	return nil
}
