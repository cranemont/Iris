package rabbitmq

import (
	"fmt"
	"log"

	"github.com/cranemont/judge-manager/constants"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	ChannelOpen(name string) error
	ExchangeDeclare(channelName, name string, typeStr string) error
	QueueDeclare(channelName, name string) error
	QueueBind(channelName, queueName string, bindingKey string, exchangeName string) error
	Subscribe(channelName string, queueName string) (<-chan amqp.Delivery, error)
	CleanUp() error
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

func (c *consumer) ChannelOpen(name string) error {
	if _, exist := c.channels[name]; exist {
		return fmt.Errorf("consumer: channel open failed: channel name already exists")
	}

	channel, err := c.connection.Channel()
	if err != nil {
		return fmt.Errorf("consumer: channel open failed: %w", err)
	}
	c.channels[name] = channel
	return nil
}

func (c *consumer) ExchangeDeclare(channelName string, exchangeName string, typeName string) error {
	channel, exist := c.channels[channelName]
	if !exist {
		return fmt.Errorf("consumer: ExchangeDeclare: channel does not exist")
	}

	if err := channel.ExchangeDeclare(
		exchangeName, // name of the exchange
		typeName,     // type
		true,         // durable
		false,        // delete when complete
		false,        // internal(deprecated)
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}
	return nil
}

func (c *consumer) QueueDeclare(channelName string, queueName string) error {
	channel, exist := c.channels[channelName]
	if !exist {
		return fmt.Errorf("consumer: QueueDeclare: channel does not exist")
	}

	_, err := channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}
	return nil
}

func (c *consumer) QueueBind(channelName string, queueName string, bindingKey string, exchangeName string) error {
	channel, exist := c.channels[channelName]
	if !exist {
		return fmt.Errorf("consumer: QueueBind: channel does not exist")
	}

	if err := channel.QueueBind(
		queueName,    // name of the queue
		bindingKey,   // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}
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
		return nil, fmt.Errorf("Queue Consume: %s", err)
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
	// FIXME: memory leak
	return <-c.Done
}
