package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/cranemont/judge-manager/ingress/mq/rabbitmq"
	"github.com/cranemont/judge-manager/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Ingress interface {
	Activate()
}

type ingress struct {
	consumer   rabbitmq.Consumer
	producer   rabbitmq.Producer
	controller RmqController
	Done       chan error
	logger     *logger.Logger
}

func NewIngress(
	consumer rabbitmq.Consumer,
	producer rabbitmq.Producer,
	controller RmqController,
	logger *logger.Logger,
) *ingress {
	return &ingress{consumer, producer, controller, make(chan error), logger}
}

func (i *ingress) Activate() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		i.consumer.CleanUp()
		i.producer.CleanUp()
	}()

	channelName := "submission"
	queueName := "submission-queue"

	err := i.consumer.OpenChannel(channelName)
	if err != nil {
		panic(err)
	}

	err = i.producer.OpenChannel()
	if err != nil {
		panic(err)
	}

	messageCh, err := i.consumer.Subscribe(channelName, queueName)
	if err != nil {
		panic(err)
	}

	// [mq.ingress]     consume -> handle -> 														  produce
	// [mq.controller]							| controller -> 	controller(result) -> |
	// [handler]													  | handler -> |
	// i.consume(messageCh, i.Done)
	for message := range messageCh {
		go i.handle(message, ctx)
	}
	// running until Consumer is done
	// <-i.Done

	// if err := i.consumer.CleanUp(); err != nil {
	// 	i.logger.Error(fmt.Sprintf("failed to clean up the consumer: %s", err))
	// }
}

func (i *ingress) handle(message amqp.Delivery, ctx context.Context) {
	result := i.controller.Call(Judge, message.Body)

	if err := i.producer.Publish(result, ctx); err != nil {
		i.logger.Error(fmt.Sprintf("failed to publish result: %s: %s", string(result), err))
	}

	if err := message.Ack(false); err != nil {
		i.logger.Error(fmt.Sprintf("failed to ack message: %s: %s", string(message.Body), err))
	}
}
