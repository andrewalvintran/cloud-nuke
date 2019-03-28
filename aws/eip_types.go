package aws

import (
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// EBSVolumes - represents all ebs volumes
type EIPAddresses struct {
	AllocationIds []string
}

// GetAllResources - Gets all the EIP addresses as a AwsResource
func (address EIPAddresses) GetAllResources(session *session.Session, region string, excludeAfter time.Time) (AwsResources, error) {
	allocationIds, err := getAllEIPAddresses(session, region, excludeAfter)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	eipAddresses := EIPAddresses{
		AllocationIds: awsgo.StringValueSlice(allocationIds),
	}

	return eipAddresses, nil
}

// ResourceName - the simple name of the aws resource
func (address EIPAddresses) ResourceName() string {
	return "eip"
}

// ResourceIdentifiers - The instance ids of the eip addresses
func (address EIPAddresses) ResourceIdentifiers() []string {
	return address.AllocationIds
}

func (address EIPAddresses) MaxBatchSize() int {
	// Tentative batch size to ensure AWS doesn't throttle
	return 200
}

// Nuke - nuke 'em all!!!
func (address EIPAddresses) Nuke(session *session.Session, identifiers []string) error {
	if err := nukeAllEIPAddresses(session, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
