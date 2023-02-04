package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	OpenChannel() error
	Subscribe() (<-chan amqp.Delivery, error)
	CleanUp() error
	// Ack(channelName string, tag uint64) error
}

type consumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
	tag        string
	Done       chan error
}

func NewConsumer(amqpURI, connectionName, queueName, ctag string) (*consumer, error) {

	// Create New RabbitMQ Connection (go <-> RabbitMQ)
	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName(connectionName)
	connection, err := amqp.DialConfig(amqpURI, config)
	if err != nil {
		return nil, fmt.Errorf("consumer: dial failed: %w", err)
	}

	return &consumer{
		connection: connection,
		channel:    nil,
		queueName:  queueName,
		tag:        ctag,
		Done:       make(chan error),
	}, nil
}

func (c *consumer) OpenChannel() error {
	var err error

	if c.channel, err = c.connection.Channel(); err != nil {
		return fmt.Errorf("channel: %s", err)
	}
	// Set prefetchCount for consume channel
	if err = c.channel.Qos(
		1,     // prefetchCount
		0,     // prefetchSize
		false, // global
	); err != nil {
		return fmt.Errorf("qos set: %s", err)
	}
	return nil
}

func (c *consumer) Subscribe() (<-chan amqp.Delivery, error) {

	// Subscribe queue for consume messages
	// Return `<- chan Delivery`
	messages, err := c.channel.Consume(
		c.queueName, // queue name
		c.tag,       // consumer
		false,       // autoAck
		false,       // exclusive
		false,       // noLocal
		false,       // noWait
		nil,         // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue consume: %s", err)
	}
	return messages, nil
}

func (c *consumer) CleanUp() error {
	// Close channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %w", err)
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
