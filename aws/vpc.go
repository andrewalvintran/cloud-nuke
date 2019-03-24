package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/cloud-nuke/logging"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

// Returns a formatted string of VPC ids
func getAllVPCs(session *session.Session, region string, excludeAfter time.Time) ([]*string, error) {
	svc := ec2.New(session)
	output, err := svc.DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	var vpcIds []*string
	for _, vpc := range output.Vpcs {
		vpcIds = append(vpcIds, vpc.VpcId)
	}

	return vpcIds, nil
}

// Deletes all VPCs
func nukeAllVPCs(session *session.Session, vpcIds []*string) error {
	svc := ec2.New(session)

	if len(vpcIds) == 0 {
		logging.Logger.Infof("No VPCs to nuke in region %s", *session.Config.Region)
		return nil
	}

	logging.Logger.Infof("Deleting all the VPCs in region %s", *session.Config.Region)
	var deletedVPCIds []*string

	for _, vpcId := range vpcIds {
		params := &ec2.DeleteVpcInput{
			VpcId: vpcId,
		}

		_, err := svc.DeleteVpc(params)

		if err != nil {
			logging.Logger.Errorf("[Failed] %s", err)
		} else {
			deletedVPCIds = append(deletedVPCIds, vpcId)
			logging.Logger.Infof("Deleted VPC id: %s", *vpcId)
		}
	}

	logging.Logger.Infof("[OK] %d VPC(s) deleted in %s", len(deletedVPCIds), *session.Config.Region)
	return nil
}
