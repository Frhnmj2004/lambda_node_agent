package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

// Client defines the interface for NATS messaging operations.
type Client interface {
	// SubscribeToJobs subscribes to the agent's job subject and invokes the handler for each job message.
	SubscribeToJobs(ctx context.Context, subject string, handler func([]byte)) error
	// PublishStatusUpdate publishes a status update message to the status subject.
	PublishStatusUpdate(ctx context.Context, status []byte) error
}

// natsClient implements Client using the nats.go library.
type natsClient struct {
	nc *nats.Conn
}

// NewNatsClient creates a new natsClient instance.
func NewNatsClient(url string) (Client, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &natsClient{nc: nc}, nil
}

// SubscribeToJobs subscribes to the agent's job subject and invokes the handler for each job message.
func (n *natsClient) SubscribeToJobs(ctx context.Context, subject string, handler func([]byte)) error {
	sub, err := n.nc.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", subject, err)
	}
	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
	}()
	return nil
}

// PublishStatusUpdate publishes a status update message to the status subject.
func (n *natsClient) PublishStatusUpdate(ctx context.Context, status []byte) error {
	subject := "jobs.status.updates"
	return n.nc.Publish(subject, status)
}
