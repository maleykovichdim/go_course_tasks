package kafka_client

import (
	"context"
	"errors"

	"github.com/segmentio/kafka-go"
)

// Client
type Client struct {
	Reader *kafka.Reader
	Writer *kafka.Writer
}

// New
func New(brokers []string, topic string, groupId string) (*Client, error) {
	if len(brokers) == 0 || brokers[0] == "" || topic == "" || groupId == "" {
		return nil, errors.New("set params for Kafka connection")
	}

	c := Client{}

	c.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupId,
		MinBytes: 10e1,
		MaxBytes: 10e6,
	})

	c.Writer = &kafka.Writer{
		Addr:                   kafka.TCP(brokers[0]),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &c, nil
}

func (c *Client) Write(ctx context.Context, m string) error {
	msg := kafka.Message{
		Value: []byte(m),
	}
	return c.Writer.WriteMessages(ctx, msg)
}

func (c *Client) Read(ctx context.Context) (kafka.Message, error) {
	return c.Reader.FetchMessage(ctx)
}

func (c *Client) CommitMessage(ctx context.Context, m kafka.Message) error {
	return c.Reader.CommitMessages(ctx, m)
}

func (c *Client) ReadConfirmed(ctx context.Context) (string, error) {

	msg, err := c.Reader.FetchMessage(ctx)
	if err != nil {
		return "", err
	}
	err = c.Reader.CommitMessages(ctx, msg)
	return string(msg.Value), err
}
