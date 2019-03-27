package aws

// we have a subnetid
import (
	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// Subnets - represents all subnets
type Subnets struct {
	SubnetIds []string
}

// ResourceName - the simple name of the aws resource
func (subnet Subnets) ResourceName() string {
	return "subnet"
}

// ResourceIdentifiers - The subnet ids
func (subnet Subnets) ResourceIdentifiers() []string {
	return subnet.SubnetIds
}

// MaxBatchSize - Tentative batch size to ensure AWS doesn't throttle
func (subnet Subnets) MaxBatchSize() int {
	return 200
}

// Nuke - nuke 'em all!!!
func (subnet Subnets) Nuke(session *session.Session, identifiers []string) error {
	if err := nukeAllSubnets(session, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
