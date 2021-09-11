package instance

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type StructMockClient struct{}

func (s *StructMockClient) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{
						Tags: []types.Tag{
							{
								Key:   aws.String("Name"),
								Value: aws.String("EC2InstanceName"),
							},
						},
					},
					{},
				},
			},
			{},
		},
	}, nil
}

func TestGetNameFromTag(t *testing.T) {
	mockClient := &StructMockClient{}
	commandInvocation := ssmTypes.CommandInvocation{
		CommandId:  aws.String("foo"),
		InstanceId: aws.String("bar"),
	}
	name, err := GetNameFromTag(mockClient, commandInvocation)
	if err != nil {
		t.Errorf("%v", err)
	}
	if *name != "EC2InstanceName" {
		t.Errorf("%v", name)
	}
}
