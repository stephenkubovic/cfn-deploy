package stackevents

import (
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// Event is an internal representation of a CloudFormation stack event
type Event struct {
	Timestamp            time.Time
	LogicalResourceID    string
	ResourceStatus       string
	ResourceType         string
	ResourceStatusReason string
	PhysicalResourceID   string
}

type cloudformationClient interface {
	DescribeStackEvents(*cloudformation.DescribeStackEventsInput) (*cloudformation.DescribeStackEventsOutput, error)
}

// New creates an internal stack event from a cloudformation.StackEvent
func New(source *cloudformation.StackEvent) (e Event) {
	e.Timestamp = *source.Timestamp
	e.LogicalResourceID = *source.LogicalResourceId
	e.ResourceStatus = *source.ResourceStatus
	e.ResourceType = *source.ResourceType
	e.PhysicalResourceID = *source.PhysicalResourceId

	if source.ResourceStatusReason == nil {
		e.ResourceStatusReason = ""
	} else {
		e.ResourceStatusReason = *source.ResourceStatusReason
	}

	return
}

// Read reads stack events for `stack` whose timestamp is after `t`
func Read(cfn cloudformationClient, t time.Time, stack string) ([]Event, error) {
	output, err := cfn.DescribeStackEvents(&cloudformation.DescribeStackEventsInput{
		StackName: &stack,
	})
	if err != nil {
		return nil, err
	}

	events := make([]Event, 0)

	for i := len(output.StackEvents) - 1; i >= 0; i-- {
		if output.StackEvents[i].Timestamp.After(t) {
			events = append(events, New(output.StackEvents[i]))
		}
	}

	return events, nil
}

func (e *Event) IsOk() bool {
	return StatusType(e.ResourceStatus) == Ok
}

func (e *Event) IsFailure() bool {
	return StatusType(e.ResourceStatus) == Fail
}

func (e *Event) IsProgress() bool {
	return StatusType(e.ResourceStatus) == Progress
}
