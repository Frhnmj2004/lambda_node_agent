package nats

import (
	"context"
	"testing"
)

func TestNatsClient_SubscribeAndPublish(t *testing.T) {
	// TODO: Start a local NATS server or use a mock
	client := &natsClient{}

	ctx := context.Background()
	subject := "jobs.dispatch.testagent"

	jobReceived := false
	err := client.SubscribeToJobs(ctx, subject, func(msg []byte) {
		jobReceived = true
	})
	if err != nil {
		t.Errorf("SubscribeToJobs failed: %v", err)
	}

	err = client.PublishStatusUpdate(ctx, []byte(`{"status":"success"}`))
	if err != nil {
		t.Errorf("PublishStatusUpdate failed: %v", err)
	}

	_ = jobReceived // avoid linter error for unused variable
	// TODO: Assert jobReceived and message published
}
