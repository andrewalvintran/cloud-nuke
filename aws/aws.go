package aws

import (
	"math/rand"
	"strings"
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/cloud-nuke/logging"
	"github.com/gruntwork-io/gruntwork-cli/collections"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// GetAllRegions - Returns a list of all AWS regions
func GetAllRegions() []string {
	// chinese and government regions are not accessible with regular accounts
	reservedRegions := []string{
		"cn-north-1", "cn-northwest-1", "us-gov-west-1", "us-gov-east-1",
	}

	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()

	var regions []string
	for _, p := range partitions {
		for id := range p.Regions() {
			if !collections.ListContainsElement(reservedRegions, id) {
				regions = append(regions, id)
			}
		}
	}

	return regions
}

func getRandomRegion() string {
	allRegions := GetAllRegions()
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(allRegions))
	return allRegions[randIndex]
}

func split(identifiers []string, limit int) [][]string {
	if limit < 0 {
		limit = -1 * limit
	} else if limit == 0 {
		return [][]string{identifiers}
	}

	var chunk []string
	chunks := make([][]string, 0, len(identifiers)/limit+1)
	for len(identifiers) >= limit {
		chunk, identifiers = identifiers[:limit], identifiers[limit:]
		chunks = append(chunks, chunk)
	}
	if len(identifiers) > 0 {
		chunks = append(chunks, identifiers[:len(identifiers)])
	}

	return chunks
}

// GetAllResources - Lists all aws resources
func GetAllResources(regions []string, excludedRegions []string, excludeAfter time.Time) (*AwsAccountResources, error) {
	account := AwsAccountResources{
		Resources: make(map[string]AwsRegionResource),
	}

	// The order in which resources are nuked is important
	// because of dependencies between resources
	allAWSResources := []AwsResources{
		ASGroups{},
		LaunchConfigs{},
		LoadBalancers{},
		LoadBalancersV2{},
		EC2Instances{},
		EBSVolumes{},
		EIPAddresses{},
		AMIs{},
		Snapshots{},
		ECSServices{},
		EKSClusters{},
	}

	for _, region := range regions {
		// Ignore all cli excluded regions
		if collections.ListContainsElement(excludedRegions, region) {
			logging.Logger.Infoln("Skipping region: " + region)
			continue
		}

		session, err := session.NewSession(&awsgo.Config{
			Region: awsgo.String(region)},
		)

		if err != nil {
			return nil, errors.WithStackTrace(err)
		}

		resourcesInRegion := AwsRegionResource{}

		for _, awsResource := range allAWSResources {
			resources, err := awsResource.GetAllResources(session, region, excludeAfter)

			if err != nil {
				return nil, errors.WithStackTrace(err)
			}

			if resources != nil {
				resourcesInRegion.Resources = append(resourcesInRegion.Resources, resources)
			}
		}

		account.Resources[region] = resourcesInRegion
	}

	return &account, nil
}

// NukeAllResources - Nukes all aws resources
func NukeAllResources(account *AwsAccountResources, regions []string) error {
	for _, region := range regions {
		session, err := session.NewSession(&awsgo.Config{
			Region: awsgo.String(region)},
		)

		if err != nil {
			return errors.WithStackTrace(err)
		}

		resourcesInRegion := account.Resources[region]
		for _, resources := range resourcesInRegion.Resources {
			length := len(resources.ResourceIdentifiers())

			// Split api calls into batches
			logging.Logger.Infof("Terminating %d resources in batches", length)
			batches := split(resources.ResourceIdentifiers(), resources.MaxBatchSize())

			for i := 0; i < len(batches); i++ {
				batch := batches[i]
				if err := resources.Nuke(session, batch); err != nil {
					// TODO: Figure out actual error type
					if strings.Contains(err.Error(), "RequestLimitExceeded") {
						logging.Logger.Info("Request limit reached. Waiting 1 minute before making new requests")
						time.Sleep(1 * time.Minute)
						continue
					}

					return errors.WithStackTrace(err)
				}

				if i != len(batches)-1 {
					logging.Logger.Info("Sleeping for 10 seconds before processing next batch...")
					time.Sleep(10 * time.Second)
				}
			}
		}
	}

	return nil
}
