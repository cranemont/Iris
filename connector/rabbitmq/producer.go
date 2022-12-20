package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/cranemont/judge-manager/common/constants"
	"github.com/cranemont/judge-manager/service/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer interface {
	OpenChannel() error
	Publish([]byte, context.Context) error
	CleanUp() error
}

type producer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	Done       chan error
	publishes  chan uint64
	logger     *logger.Logger
}

func NewProducer(amqpURI string, connectionName string, logger *logger.Logger) (*producer, error) {

	// Create New RabbitMQ Connection (go <-> RabbitMQ)
	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName(connectionName)
	connection, err := amqp.DialConfig(amqpURI, config)
	if err != nil {
		return nil, fmt.Errorf("consumer: dial failed: %w", err)
	}

	return &producer{
		connection: connection,
		channel:    nil,
		Done:       make(chan error),
		publishes:  make(chan uint64, 8),
		logger:     logger,
	}, nil
}

func (p *producer) OpenChannel() error {
	var err error

	if p.channel, err = p.connection.Channel(); err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	// put this channel into confirm mode
	// client can ensure all messages successfully received by server
	if err := p.channel.Confirm(false); err != nil {
		return fmt.Errorf("channel could not be put into confirm mode: %s", err)
	}
	// add listner for confirmation
	confirms := p.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	go p.confirmHandler(confirms)

	return nil
}

func (p *producer) confirmHandler(confirms chan amqp.Confirmation) {
	m := make(map[uint64]bool)
	for {
		select {
		case <-p.Done:
			p.logger.Info("confirmHandler is stopping")
			return
		case publishSeqNo := <-p.publishes:
			// log.Printf("waiting for confirmation of %d", publishSeqNo)
			m[publishSeqNo] = false
		case confirmed := <-confirms:
			if confirmed.DeliveryTag > 0 {
				if confirmed.Ack {
					p.logger.Debug(fmt.Sprintf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag))
				} else {
					p.logger.Error(fmt.Sprintf("failed delivery of delivery tag: %d", confirmed.DeliveryTag))
				}
				delete(m, confirmed.DeliveryTag)
			}
		}
	}
}

func (p *producer) Publish(result []byte, ctx context.Context) error {

	seqNo := p.channel.GetNextPublishSeqNo()
	log.Printf("publishing %dB body (%q)", len(result), result)

	if err := p.channel.PublishWithContext(ctx,
		constants.EXCHANGE,   // publish to an exchange
		constants.RESULT_KEY, // routing to 0 or more queues
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            result,
			DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
			Priority:        0,               // 0-9
		},
	); err != nil {
		return fmt.Errorf("exchange publish: %s", err)
	}

	p.logger.Debug(fmt.Sprintf("published %dB OK", len(result)))
	p.publishes <- seqNo

	return nil
}

func (p *producer) CleanUp() error {
	if err := p.channel.Close(); err != nil {
		return fmt.Errorf("channel close failed: %s", err)
	}

	if err := p.connection.Close(); err != nil {
		return fmt.Errorf("connection close error: %s", err)
	}

	return <-p.Done
}
