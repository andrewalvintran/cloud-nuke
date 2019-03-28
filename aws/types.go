package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
)

type AwsAccountResources struct {
	Resources map[string]AwsRegionResource
}

type AwsResources interface {
	GetAllResources(session *session.Session, region string, excludeAfter time.Time) (AwsResources, error)
	ResourceName() string
	ResourceIdentifiers() []string
	MaxBatchSize() int
	Nuke(session *session.Session, identifiers []string) error
}

type AwsRegionResource struct {
	Resources []AwsResources
}
