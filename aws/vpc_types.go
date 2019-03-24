package aws

import (
	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// VPCs - represent all the non default VPCs
type VPCs struct {
	VPCIds []string
}

// ResourceName - the simple name of the aws resource
func (vpc VPCs) ResourceName() string {
	return "vpc"
}

// ResourceIdentifiers - The ids of the VPCs
func (vpc VPCs) ResourceIdentifiers() []string {
	return vpc.VPCIds
}

// MaxBatchSize - Tentative batch size to ensure AWS doesn't throttle
func (vpc VPCs) MaxBatchSize() int {
	return 200
}

// Nuke - nuke 'em all!!!
func (vpc VPCs) Nuke(session *session.Session, identifiers []string) error {
	if err := nukeAllVPCs(session, awsgo.StringSlice(identifiers)); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
