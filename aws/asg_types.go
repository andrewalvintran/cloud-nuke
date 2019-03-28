package aws

import (
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// ASGroups - represents all auto scaling groups
type ASGroups struct {
	GroupNames []string
}

// GetResources - Gets all the asg group names as a AwsResource
func (group ASGroups) GetResources(session *session.Session, region string, excludeAfter time.Time) (AwsResources, error) {
	groupNames, err := getAllAutoScalingGroups(session, region, excludeAfter)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	asGroups := ASGroups{
		GroupNames: awsgo.StringValueSlice(groupNames),
	}

	return asGroups, nil
}

// ResourceName - the simple name of the aws resour    ce
func (group ASGroups) ResourceName() string {
	return "asg"
}

func (group ASGroups) MaxBatchSize() int {
	// Tentative batch size to ensure AWS doesn't throttle
	return 200
}

// ResourceIdentifiers - The group names of the auto scaling groups
func (group ASGroups) ResourceIdentifiers() []string {
	return group.GroupNames
}

// Nuke - nuke 'em all!!!
func (group ASGroups) Nuke(session *session.Session, identifiers []string) error {
	if err := nukeAllAutoScalingGroups(session, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
