package logs

import (
	"context"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type GetLogEventsAPI interface {
	GetLogEvents(ctx context.Context, params *cloudwatchlogs.GetLogEventsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetLogEventsOutput, error)
}

func GetAllLogs(client GetLogEventsAPI, logGroupName string, logStreamName string) ([]types.OutputLogEvent, error) {
	events := []types.OutputLogEvent{}
	var nextToken *string = nil
	var lastToken *string
	for ok := true; ok; ok = !reflect.DeepEqual(nextToken, lastToken) {
		output, err := client.GetLogEvents(context.TODO(), createGetLogEventsInput(logGroupName, logStreamName, nextToken))
		if err != nil {
			return events, err
		}
		events = append(events, output.Events...)
		lastToken = nextToken
		nextToken = output.NextForwardToken
	}
	return events, nil
}

func createGetLogEventsInput(logGroupName string, logStreamName string, nextToken *string) *cloudwatchlogs.GetLogEventsInput {
	return &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
		StartFromHead: aws.Bool(true),
		NextToken:     nextToken,
	}
}
