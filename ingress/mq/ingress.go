package mq

import (
	"fmt"

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
	logging    *logger.Logger
}

func NewIngress(
	consumer rabbitmq.Consumer,
	producer rabbitmq.Producer,
	controller RmqController,
	logging *logger.Logger,
) *ingress {
	return &ingress{consumer, producer, controller, make(chan error), logging}
}

func (i *ingress) Activate() {
	// goroutine 등록
	// 여기서 자동으로 consume해서 시작됨(go controller.call)

	channelName := "submission"
	queueName := "submission-queue"
	exchangeName := "submission-exchange"
	exchangeType := "direct"
	bindingKey := "submission"

	err := i.consumer.ChannelOpen(channelName)
	if err != nil {
		panic(err)
	}
	err = i.consumer.ExchangeDeclare(channelName, exchangeName, exchangeType)
	if err != nil {
		panic(err)
	}
	err = i.consumer.QueueDeclare(channelName, queueName)
	if err != nil {
		panic(err)
	}
	err = i.consumer.QueueBind(channelName, queueName, bindingKey, exchangeName)
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
	i.consume(messageCh, i.Done)

	// running until Consumer is done
	<-i.Done

	if err := i.consumer.CleanUp(); err != nil {
		i.logging.Error(fmt.Sprintf("failed to clean up the consumer: %s", err))
	}
}

func (i *ingress) consume(messages <-chan amqp.Delivery, done chan error) {
	clean := func() {
		done <- nil
	}
	defer clean()

	for message := range messages {
		go i.handle(message)
	}
}

func (i *ingress) handle(message amqp.Delivery) {
	result := i.controller.Call(Judge, message.Body)

	i.producer.Publish(result)
	if err := message.Ack(false); err != nil {
		i.logging.Error(fmt.Sprintf("failed to ack message: %s", err))
	}
}