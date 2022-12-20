package rabbitmq

import (
	"fmt"
	"log"

	"github.com/cranemont/judge-manager/common/constants"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	OpenChannel(name string) error
	Subscribe(channelName string, queueName string) (<-chan amqp.Delivery, error)
	CleanUp() error
	// Ack(channelName string, tag uint64) error
}

type consumer struct {
	connection *amqp.Connection
	channels   map[string](*amqp.Channel)
	tag        string
	Done       chan error
}

func NewConsumer(amqpURI string, ctag string, connectionName string) (*consumer, error) {

	// Create New RabbitMQ Connection (go <-> RabbitMQ)
	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName(connectionName)
	connection, err := amqp.DialConfig(amqpURI, config)
	if err != nil {
		return nil, fmt.Errorf("consumer: dial failed: %w", err)
	}

	return &consumer{
		connection: connection,
		channels:   make(map[string](*amqp.Channel), constants.MAX_MQ_CHANNEL),
		tag:        ctag,
		Done:       make(chan error),
	}, nil
}

func (c *consumer) OpenChannel(name string) error {
	if _, exist := c.channels[name]; exist {
		return fmt.Errorf("consumer: channel open failed: channel name already exists")
	}

	channel, err := c.connection.Channel()
	if err != nil {
		return fmt.Errorf("consumer: channel open failed: %w", err)
	}
	// Set prefetchCount for consume channel
	if err = channel.Qos(
		1,     // prefetchCount
		0,     // prefetchSize
		false, // global
	); err != nil {
		return fmt.Errorf("qos set: %s", err)
	}
	c.channels[name] = channel
	return nil
}

func (c *consumer) Subscribe(channelName string, queueName string) (<-chan amqp.Delivery, error) {
	channel, exist := c.channels[channelName]
	if !exist {
		return nil, fmt.Errorf("consumer: Consume: channel does not exist")
	}

	// Subscribe queue for consume messages
	// Return `<- chan Delivery`
	messages, err := channel.Consume(
		queueName, // queue name
		c.tag,     // consumer
		false,     // autoAck
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue consume: %s", err)
	}
	return messages, nil
}

func (c *consumer) CleanUp() error {
	// Close channel
	for name, channel := range c.channels {
		if err := channel.Cancel(c.tag, true); err != nil {
			return fmt.Errorf("Consumer cancel failed: %s: %w", name, err)
		}
	}

	// Close Connection
	if err := c.connection.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}
	defer log.Print("RabbitMQ connection clear done")

	// wait for handle() to exit
	return <-c.Done
}

// func (c *consumer) Ack(channelName string, tag uint64) error {
// 	channel, exist := c.channels[channelName]
// 	if !exist {
// 		return fmt.Errorf("consumer: Ack: channel does not exist")
// 	}
// 	err := channel.Ack(tag, false)
// 	return err
// }
