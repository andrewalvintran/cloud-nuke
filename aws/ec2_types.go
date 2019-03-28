package aws

import (
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// EC2Instances - represents all ec2 instances
type EC2Instances struct {
	InstanceIds []string
}

// GetAllResources - Gets all the ec2 instance ids as a AwsResource
func (instance EC2Instances) GetAllResources(session *session.Session, region string, excludeAfter time.Time) (AwsResources, error) {
	instanceIds, err := getAllEc2Instances(session, region, excludeAfter)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	ec2Instances := EC2Instances{
		InstanceIds: awsgo.StringValueSlice(instanceIds),
	}

	return ec2Instances, nil
}

// ResourceName - the simple name of the aws resource
func (instance EC2Instances) ResourceName() string {
	return "ec2"
}

// ResourceIdentifiers - The instance ids of the ec2 instances
func (instance EC2Instances) ResourceIdentifiers() []string {
	return instance.InstanceIds
}

func (instance EC2Instances) MaxBatchSize() int {
	// Tentative batch size to ensure AWS doesn't throttle
	return 200
}

// Nuke - nuke 'em all!!!
func (instance EC2Instances) Nuke(session *session.Session, identifiers []string) error {
	if err := nukeAllEc2Instances(session, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
