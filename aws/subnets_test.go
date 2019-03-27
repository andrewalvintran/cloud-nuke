package aws

import (
	"testing"
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/gruntwork-cli/errors"
	gruntworkerrors "github.com/gruntwork-io/gruntwork-cli/errors"
	"github.com/stretchr/testify/assert"
)

func createTestSubnet(t *testing.T, session *session.Session) (ec2.Subnet, ec2.Vpc) {
	svc := ec2.New(session)

	vpc := createTestVPC(t, session)

	input := &ec2.CreateSubnetInput{
		CidrBlock: awsgo.String("10.0.0.0/16"),
		VpcId:     awsgo.String(*vpc.VpcId),
	}

	subnet, err := svc.CreateSubnet(input)
	if err != nil {
		assert.Fail(t, gruntworkerrors.WithStackTrace(err).Error())
	}

	return *subnet.Subnet, vpc
}

func TestListSubnets(t *testing.T) {
	t.Parallel()

	region := getRandomRegion()
	session, err := session.NewSession(&awsgo.Config{
		Region: awsgo.String(region)},
	)

	if err != nil {
		assert.Fail(t, errors.WithStackTrace(err).Error())
	}

	subnet, vpc := createTestSubnet(t, session)

	// clean up after test
	defer nukeAllVPCs(session, []*string{vpc.VpcId})
	defer nukeAllSubnets(session, []*string{subnet.SubnetId})

	subnetIds, err := getAllSubnets(session, region, time.Now().Add(1*time.Hour))
	if err != nil {
		assert.Fail(t, "Unable to fetch list of subnets")
	}

	assert.Contains(t, awsgo.StringValueSlice(subnetIds), *subnet.SubnetId)
}

func TestNukeSubnets(t *testing.T) {
	t.Parallel()

	region := getRandomRegion()
	session, err := session.NewSession(&awsgo.Config{
		Region: awsgo.String(region)},
	)

	if err != nil {
		assert.Fail(t, errors.WithStackTrace(err).Error())
	}

	subnet, vpc := createTestSubnet(t, session)

	defer nukeAllVPCs(session, []*string{vpc.VpcId})

	if err := nukeAllSubnets(session, []*string{subnet.SubnetId}); err != nil {
		assert.Fail(t, gruntworkerrors.WithStackTrace(err).Error())
	}

	subnetIds, err := getAllSubnets(session, region, time.Now().Add(1*time.Hour))

	if err != nil {
		assert.Fail(t, "Unable to fetch list of subnet ids")
	}

	assert.NotContains(t, awsgo.StringValueSlice(subnetIds), *subnet.SubnetId)
}
