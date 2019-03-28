package aws

import (
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// LaunchConfigs - represents all launch configurations
type LaunchConfigs struct {
	LaunchConfigurationNames []string
}

// GetResources - Gets all the launch configs as a AwsResource
func (config LaunchConfigs) GetResources(session *session.Session, region string, excludeAfter time.Time) (AwsResources, error) {
	configNames, err := getAllLaunchConfigurations(session, region, excludeAfter)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	configs := LaunchConfigs{
		LaunchConfigurationNames: awsgo.StringValueSlice(configNames),
	}

	return configs, nil
}

// ResourceName - the simple name of the aws resource
func (config LaunchConfigs) ResourceName() string {
	return "lc"
}

func (config LaunchConfigs) MaxBatchSize() int {
	// Tentative batch size to ensure AWS doesn't throttle
	return 200
}

// ResourceIdentifiers - The names of the launch configurations
func (config LaunchConfigs) ResourceIdentifiers() []string {
	return config.LaunchConfigurationNames
}

// Nuke - nuke 'em all!!!
func (config LaunchConfigs) Nuke(session *session.Session, identifiers []string) error {
	if err := nukeAllLaunchConfigurations(session, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
