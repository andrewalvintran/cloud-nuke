package aws

import (
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// AMIs - represents all user owned AMIs
type AMIs struct {
	ImageIds []string
}

// GetAllResources - Gets all the ami image ids as a AwsResource
func (image AMIs) GetAllResources(session *session.Session, region string, excludeAfter time.Time) (AwsResources, error) {
	imageIds, err := getAllAMIs(session, region, excludeAfter)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	amis := AMIs{
		ImageIds: awsgo.StringValueSlice(imageIds),
	}

	return amis, nil
}

// ResourceName - the simple name of the aws resource
func (image AMIs) ResourceName() string {
	return "ami"
}

// ResourceIdentifiers - The AMI image ids
func (image AMIs) ResourceIdentifiers() []string {
	return image.ImageIds
}

func (image AMIs) MaxBatchSize() int {
	// Tentative batch size to ensure AWS doesn't throttle
	return 200
}

// Nuke - nuke 'em all!!!
func (image AMIs) Nuke(session *session.Session, identifiers []string) error {
	if err := nukeAllAMIs(session, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}

type ImageAvailableError struct{}

func (e ImageAvailableError) Error() string {
	return "Image didn't become available within wait attempts"
}
