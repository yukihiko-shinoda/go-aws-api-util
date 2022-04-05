package command

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type StructMockClient struct {
	TimeListCommandInvocations int
}

func (s *StructMockClient) ListCommands(ctx context.Context, params *ssm.ListCommandsInput, optFns ...func(*ssm.Options)) (*ssm.ListCommandsOutput, error) {
	return &ssm.ListCommandsOutput{
		Commands: []types.Command{
			{
				CommandId:         aws.String(""),
				RequestedDateTime: &time.Time{},
			},
		},
	}, nil
}

func (s *StructMockClient) ListCommandInvocations(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (*ssm.ListCommandInvocationsOutput, error) {
	s.TimeListCommandInvocations++
	if s.TimeListCommandInvocations <= 1 {
		return &ssm.ListCommandInvocationsOutput{
			CommandInvocations: []types.CommandInvocation{
				{
					InstanceId:        aws.String("i-1234567890abcdef0"),
					InstanceName:      aws.String(""),
					RequestedDateTime: &time.Time{},
				},
			},
			NextToken: aws.String(""),
		}, nil
	}
	return &ssm.ListCommandInvocationsOutput{
		CommandInvocations: []types.CommandInvocation{},
		NextToken:          nil,
	}, nil
}

func newMockClient() *StructMockClient {
	return &StructMockClient{
		TimeListCommandInvocations: 0,
	}
}
func TestGetLatestApplyAnsiblePlaybooksInvocations(t *testing.T) {
	mockClient := newMockClient()
	expected := []types.CommandInvocation{
		{
			InstanceId:        aws.String("i-1234567890abcdef0"),
			InstanceName:      aws.String(""),
			RequestedDateTime: &time.Time{},
		},
		{},
	}
	invocations, err := GetApplyAnsiblePlaybooksInvocations(mockClient, nil)
	if err != nil {
		t.Errorf("%v", invocations)
	}
	if reflect.DeepEqual(invocations, expected) {
		t.Errorf("%v", invocations)
	}
}

func TestBuildLogStreamNameRunShellScriptStdout(t *testing.T) {
	expect := "foo/bar/runShellScript/stdout"
	commandInvocation := types.CommandInvocation{
		CommandId:  aws.String("foo"),
		InstanceId: aws.String("bar"),
	}
	actual := BuildLogStreamNameRunShellScriptStdout(commandInvocation)
	if actual != expect {
		t.Errorf("%v", actual)
	}
}
