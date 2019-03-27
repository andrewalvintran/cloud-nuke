package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/cloud-nuke/logging"
	"github.com/gruntwork-io/gruntwork-cli/errors"
)

func getAllSubnets(session *session.Session, region string, excludeAfter time.Time) ([]*string, error) {
	svc := ec2.New(session)

	output, err := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{})
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	var subnetIds []*string
	for _, subnet := range output.Subnets {
		subnetIds = append(subnetIds, subnet.SubnetId)
	}

	return subnetIds, nil
}

func nukeAllSubnets(session *session.Session, subnetIds []*string) error {
	svc := ec2.New(session)

	if len(subnetIds) == 0 {
		logging.Logger.Infof("No subnets to nuke in region %s", *session.Config.Region)
		return nil
	}

	logging.Logger.Infof("Deleting all the subnets in region %s", *session.Config.Region)
	var deletedSubnetIds []*string

	for _, subnetID := range subnetIds {
		params := &ec2.DeleteSubnetInput{
			SubnetId: subnetID,
		}
		_, err := svc.DeleteSubnet(params)

		if err != nil {
			logging.Logger.Errorf("[Failed] %s", err)
		} else {
			deletedSubnetIds = append(deletedSubnetIds, subnetID)
			logging.Logger.Infof("Deleted subnet: %s", *subnetID)
		}
	}

	logging.Logger.Infof("[OK] %d subnet(s) deleted in %s", len(deletedSubnetIds), *session.Config.Region)
	return nil
}
