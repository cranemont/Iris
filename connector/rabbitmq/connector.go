package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/cranemont/judge-manager/router"
	"github.com/cranemont/judge-manager/service/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type connector struct {
	consumer Consumer
	producer Producer
	router   router.Router
	Done     chan error
	logger   *logger.Logger
}

func NewConnector(
	consumer Consumer,
	producer Producer,
	router router.Router,
	logger *logger.Logger,
) *connector {
	return &connector{consumer, producer, router, make(chan error), logger}
}

func (c *connector) Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		c.consumer.CleanUp()
		c.producer.CleanUp()
	}()

	channelName := "submission"
	queueName := "submission-queue"

	err := c.consumer.OpenChannel(channelName)
	if err != nil {
		panic(err)
	}

	err = c.producer.OpenChannel()
	if err != nil {
		panic(err)
	}

	messageCh, err := c.consumer.Subscribe(channelName, queueName)
	if err != nil {
		panic(err)
	}

	// [mq.ingress]     consume -> handle -> 														  produce
	// [mq.controller]							| controller -> 	controller(result) -> |
	// [handler]													  | handler -> |
	// i.consume(messageCh, i.Done)
	for message := range messageCh {
		go c.handle(message, ctx)
	}
	// running until Consumer is done
	// <-i.Done

	// if err := i.consumer.CleanUp(); err != nil {
	// 	i.logger.Error(fmt.Sprintf("failed to clean up the consumer: %s", err))
	// }
}

func (c *connector) Disconnect() {}

func (c *connector) handle(message amqp.Delivery, ctx context.Context) {
	result := c.router.Route(router.To(message.Type), message.MessageId, message.Body)

	if err := c.producer.Publish(result, ctx); err != nil {
		c.logger.Error(fmt.Sprintf("failed to publish result: %s: %s", string(result), err))
		// nack
	}

	if err := message.Ack(false); err != nil {
		c.logger.Error(fmt.Sprintf("failed to ack message: %s: %s", string(message.Body), err))
		// retry
	}
}
