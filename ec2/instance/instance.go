package instance

import (
	"context"
	"errors"

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
	if len(result.Reservations) == 0 {
		return nil, errors.New("Reservations no longer exist.")
	}
	reservation := result.Reservations[0]
	if len(reservation.Instances) == 0 {
		return nil, errors.New("Instances no longer exist.")
	}
	instance := reservation.Instances[0]
	return instance.Tags, err
}
