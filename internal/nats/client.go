package nats

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

// Client defines the interface for NATS operations
type Client interface {
	SubscribeToJobs(ctx context.Context, subject string, handler func(msg []byte)) error
	PublishStatusUpdate(ctx context.Context, status []byte) error
	Close()
}

// natsClient implements Client using NATS
type natsClient struct {
	conn *nats.Conn
}

// NewNatsClient creates a new NATS client
func NewNatsClient(url string) (Client, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &natsClient{
		conn: conn,
	}, nil
}

// SubscribeToJobs subscribes to a NATS subject for job messages
func (n *natsClient) SubscribeToJobs(ctx context.Context, subject string, handler func(msg []byte)) error {
	subscription, err := n.conn.Subscribe(subject, func(msg *nats.Msg) {
		log.Printf("Received job message on subject: %s", subject)
		handler(msg.Data)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject %s: %w", subject, err)
	}

	// Handle context cancellation
	go func() {
		<-ctx.Done()
		subscription.Unsubscribe()
		log.Printf("Unsubscribed from subject: %s", subject)
	}()

	return nil
}

// PublishStatusUpdate publishes a status update to NATS
func (n *natsClient) PublishStatusUpdate(ctx context.Context, status []byte) error {
	// Publish to a status topic
	subject := "agent.status"
	err := n.conn.Publish(subject, status)
	if err != nil {
		return fmt.Errorf("failed to publish status update: %w", err)
	}

	log.Printf("Published status update to subject: %s", subject)
	return nil
}

// Close closes the NATS connection
func (n *natsClient) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}
