package logs

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type StructMockClient struct{}

func (s *StructMockClient) GetLogEvents(ctx context.Context, params *cloudwatchlogs.GetLogEventsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetLogEventsOutput, error) {
	return &cloudwatchlogs.GetLogEventsOutput{
		Events: []types.OutputLogEvent{
			{
				Message: aws.String("message"),
			},
		},
		NextForwardToken: aws.String("next"),
	}, nil
}

func TestGetAllLogs(t *testing.T) {
	mockClient := &StructMockClient{}
	logGroupName := "aaaa"
	logStreamName := "bbbb"
	events, err := GetAllLogs(mockClient, logGroupName, logStreamName)
	if err != nil {
		t.Errorf("%v", err)
	}
	if events == nil {
		t.Errorf("%v", events)
	}
	if !reflect.DeepEqual(events, []types.OutputLogEvent{
		{
			Message: aws.String("message"),
		},
		{
			Message: aws.String("message"),
		},
	}) {
		t.Errorf("%v", len(events))
		for _, event := range events {
			t.Errorf("%v", *event.Message)
		}
		t.Errorf("%v", events)
	}
}

func TestCreateGetLogEventsInput(t *testing.T) {
	logGroupName := "aaaa"
	logStreamName := "bbbb"
	var nextToken *string = nil
	expected := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
		StartFromHead: aws.Bool(true),
		NextToken:     nextToken,
	}
	actual := createGetLogEventsInput(logGroupName, logStreamName, nextToken)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v", actual)
	}
}
