package command

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type ListCommandsListCommandInvocationsAPI interface {
	ListCommands(ctx context.Context, params *ssm.ListCommandsInput, optFns ...func(*ssm.Options)) (*ssm.ListCommandsOutput, error)
	ListCommandInvocations(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (*ssm.ListCommandInvocationsOutput, error)
}

func GetLatestApplyAnsiblePlaybooksInvocations(client ListCommandsListCommandInvocationsAPI) ([]types.CommandInvocation, error) {
	latestCommand, err := getLatestApplyAnsiblePlaybooks(client)
	if err != nil {
		return nil, err
	}
	invocations, err := getAllInvocations(client, *latestCommand)
	if err != nil {
		return invocations, err
	}
	return invocations, nil
}

func getLatestApplyAnsiblePlaybooks(client ListCommandsListCommandInvocationsAPI) (*types.Command, error) {
	DocumentNameAnsible := "AWS-ApplyAnsiblePlaybooks"
	input := &ssm.ListCommandsInput{
		Filters:    []types.CommandFilter{{Key: "DocumentName", Value: &DocumentNameAnsible}},
		MaxResults: 1,
	}
	output, err := client.ListCommands(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	command := output.Commands[0]
	log.Printf(
		// "%s\n\t%s\n\t%s\n\t%s\n\t%s\n",
		"Execution date time: %s\tCommand ID: %s\n",
		command.RequestedDateTime.Local(),
		*command.CommandId,
	)
	return &command, nil
}

func getAllInvocations(client ListCommandsListCommandInvocationsAPI, command types.Command) ([]types.CommandInvocation, error) {
	invocations := []types.CommandInvocation{}
	var nextToken *string = nil
	for ok := true; ok; ok = nextToken != nil {
		output, err := client.ListCommandInvocations(context.TODO(), createListCommandInvocationsInput(command, nextToken))
		if err != nil {
			return invocations, err
		}
		invocations = append(invocations, output.CommandInvocations...)
		nextToken = output.NextToken
	}
	for _, i := range invocations {
		fmt.Printf(
			// "%s\n\t%s\n\t%s\n\t%s\n\t%s\n",
			"%s\t%s\t\t%s\n",
			*i.InstanceId,
			*i.InstanceName,
			i.RequestedDateTime.Local(),
		)
	}
	return invocations, nil
}

func createListCommandInvocationsInput(command types.Command, nextToken *string) *ssm.ListCommandInvocationsInput {
	return &ssm.ListCommandInvocationsInput{
		CommandId: command.CommandId,
		Details:   false,
		NextToken: nextToken,
	}
}

func BuildLogStreamNameRunShellScriptStdout(commandInvocation types.CommandInvocation) string {
	return *commandInvocation.CommandId + "/" + *commandInvocation.InstanceId + "/runShellScript/stdout"
}
