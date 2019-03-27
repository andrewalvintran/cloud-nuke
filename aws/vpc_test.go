package aws

import (
	"testing"
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	gruntworkerrors "github.com/gruntwork-io/gruntwork-cli/errors"
	"github.com/stretchr/testify/assert"
)

func createTestVPC(t *testing.T, session *session.Session) ec2.Vpc {
	svc := ec2.New(session)

	input := &ec2.CreateVpcInput{
		CidrBlock: awsgo.String("10.0.0.0/16"),
	}

	vpc, err := svc.CreateVpc(input)
	if err != nil {
		assert.Fail(t, gruntworkerrors.WithStackTrace(err).Error())
	}

	return *vpc.Vpc
}

func TestListVPCs(t *testing.T) {
	t.Parallel()

	region := getRandomRegion()
	session, err := session.NewSession(&awsgo.Config{
		Region: awsgo.String(region)},
	)

	if err != nil {
		assert.Fail(t, gruntworkerrors.WithStackTrace(err).Error())
	}

	vpc := createTestVPC(t, session)

	// clean up after test
	defer nukeAllVPCs(session, []*string{vpc.VpcId})

	vpcIds, err := getAllVPCs(session, region, time.Now().Add(1*time.Hour))
	if err != nil {
		assert.Fail(t, "Unable to fetch list of VPCs")
	}

	assert.Contains(t, awsgo.StringValueSlice(vpcIds), *vpc.VpcId)
}

func TestNukeVPCs(t *testing.T) {
	t.Parallel()

	region := getRandomRegion()
	session, err := session.NewSession(&awsgo.Config{
		Region: awsgo.String(region)},
	)

	if err != nil {
		assert.Fail(t, gruntworkerrors.WithStackTrace(err).Error())
	}

	vpc := createTestVPC(t, session)

	if err := nukeAllVPCs(session, []*string{vpc.VpcId}); err != nil {
		assert.Fail(t, gruntworkerrors.WithStackTrace(err).Error())
	}

	vpcIds, err := getAllVPCs(session, region, time.Now().Add(1*time.Hour))
	if err != nil {
		assert.Fail(t, "Unable to fetch list of VPCs")
	}

	assert.NotContains(t, awsgo.StringValueSlice(vpcIds), *vpc.VpcId)
}
