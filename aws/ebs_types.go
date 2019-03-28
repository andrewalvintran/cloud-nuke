package aws

import (
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// EBSVolumes - represents all ebs volumes
type EBSVolumes struct {
	VolumeIds []string
}

// GetResources - Gets all the ebs volumes as a AwsResource
func (volume EBSVolumes) GetResources(session *session.Session, region string, excludeAfter time.Time) (AwsResources, error) {
	volumeIds, err := getAllEbsVolumes(session, region, excludeAfter)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	ebsVolumes := EBSVolumes{
		VolumeIds: awsgo.StringValueSlice(volumeIds),
	}

	return ebsVolumes, nil
}

// ResourceName - the simple name of the aws resource
func (volume EBSVolumes) ResourceName() string {
	return "ebs"
}

// ResourceIdentifiers - The volume ids of the ebs volumes
func (volume EBSVolumes) ResourceIdentifiers() []string {
	return volume.VolumeIds
}

func (volume EBSVolumes) MaxBatchSize() int {
	// Tentative batch size to ensure AWS doesn't throttle
	return 200
}

// Nuke - nuke 'em all!!!
func (volume EBSVolumes) Nuke(session *session.Session, identifiers []string) error {
	if err := nukeAllEbsVolumes(session, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
