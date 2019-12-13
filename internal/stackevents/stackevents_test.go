package stackevents

import (
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"testing"
	"time"
)

type mockCloudFormationClient struct{}

type mockEvent struct {
	Timestamp  time.Time
	ID         string
	Status     string
	Type       string
	PhysicalID string
}

func (c mockCloudFormationClient) DescribeStackEvents(input *cfn.DescribeStackEventsInput) (*cfn.DescribeStackEventsOutput, error) {
	var events []*cfn.StackEvent

	now := mockEvent{
		Timestamp:  time.Now(),
		ID:         "1",
		Status:     cfn.StackStatusCreateComplete,
		Type:       "AWS::EC2::Instance",
		PhysicalID: "1",
	}

	hourAgo := mockEvent{
		Timestamp:  time.Now().Add(-1 * time.Hour),
		ID:         "2",
		Status:     cfn.StackStatusCreateComplete,
		Type:       "AWS::EC2::Instance",
		PhysicalID: "2",
	}

	events = append(events, &cfn.StackEvent{
		Timestamp:          &now.Timestamp,
		LogicalResourceId:  &now.ID,
		ResourceStatus:     &now.Status,
		ResourceType:       &now.Type,
		PhysicalResourceId: &now.PhysicalID,
	}, &cfn.StackEvent{
		Timestamp:          &hourAgo.Timestamp,
		LogicalResourceId:  &hourAgo.ID,
		ResourceStatus:     &hourAgo.Status,
		ResourceType:       &hourAgo.Type,
		PhysicalResourceId: &hourAgo.PhysicalID,
	})

	output := cfn.DescribeStackEventsOutput{
		StackEvents: events,
	}

	return &output, nil
}

func TestGetStackEvents(t *testing.T) {
	client := mockCloudFormationClient{}
	results, err := Read(client, time.Now().Add(-1*time.Minute), "test-stack")

	if err != nil {
		t.Error(err)
	}

	n := len(results)
	if n != 1 {
		t.Errorf("Expected a single stack event, got %d", n)
	}
}

func TestGetStackEventsEmptyResult(t *testing.T) {
	client := mockCloudFormationClient{}
	results, err := Read(client, time.Now().Add(time.Minute), "test-stack")

	if err != nil {
		t.Error(err)
	}

	n := len(results)
	if n != 0 {
		t.Errorf("Expected no stack events, got %d", n)
	}
}
