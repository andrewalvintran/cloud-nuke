package aws

import (
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// EKSClusters - Represents all EKS clusters found in a region
type EKSClusters struct {
	Clusters []string
}

// GetResources - Gets all the EKS clusters as a AwsResource
func (clusters EKSClusters) GetResources(session *session.Session, region string, excludeAfter time.Time) (AwsResources, error) {
	if eksSupportedRegion(region) {
		eksClusterNames, err := getAllEksClusters(session, excludeAfter)
		if err != nil {
			return nil, errors.WithStackTrace(err)
		}

		eksClusters := EKSClusters{
			Clusters: awsgo.StringValueSlice(eksClusterNames),
		}

		return eksClusters, nil
	}

	return nil, nil
}

// ResourceName - The simple name of the aws resource
func (clusters EKSClusters) ResourceName() string {
	return "ekscluster"
}

// ResourceIdentifiers - The Name of the collected EKS clusters
func (clusters EKSClusters) ResourceIdentifiers() []string {
	return clusters.Clusters
}

func (clusters EKSClusters) MaxBatchSize() int {
	return 200
}

// Nuke - nuke all EKS Cluster resources
func (clusters EKSClusters) Nuke(awsSession *session.Session, identifiers []string) error {
	if err := nukeAllEksClusters(awsSession, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}
	return nil
}
