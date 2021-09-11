package instance

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func GetNameFromTag(client DescribeInstancesAPI, commandInvocation ssmTypes.CommandInvocation) (*string, error) {
	tags, err := GetTags(client, commandInvocation)
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		if *tag.Key == "Name" {
			return tag.Value, nil
		}
	}
	return nil, nil
}

func GetTags(client DescribeInstancesAPI, commandInvocation ssmTypes.CommandInvocation) ([]ec2Types.Tag, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{*commandInvocation.InstanceId},
	}
	result, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	reservation := result.Reservations[0]
	instance := reservation.Instances[0]
	return instance.Tags, err
}
