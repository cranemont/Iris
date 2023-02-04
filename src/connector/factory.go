package connector

import (
	"github.com/cranemont/iris/src/connector/rabbitmq"
	"github.com/cranemont/iris/src/router"
	"github.com/cranemont/iris/src/service/logger"
	"github.com/cranemont/iris/src/utils"
)

type Providers struct {
	Router router.Router
	Logger logger.Logger
}

type Module string

const (
	RABBIT_MQ Module = "RabbitMQ"
	HTTP      Module = "Http"
	FILE      Module = "File"
	CONSOLE   Module = "Console"
)

func Factory(c Module, p Providers) Connector {
	switch c {
	case RABBIT_MQ:
		uri := "amqp://" +
			utils.Getenv("RABBITMQ_DEFAULT_USER", "skku") + ":" +
			utils.Getenv("RABBITMQ_DEFAULT_PASS", "1234") + "@" +
			utils.Getenv("RABBITMQ_HOST", "localhost") + ":" +
			utils.Getenv("RABBITMQ_PORT", "5672") + "/"
		consumer, err := rabbitmq.NewConsumer(uri, "ctag", "go-consumer")
		if err != nil {
			panic(err)
		}
		producer, err := rabbitmq.NewProducer(uri, "go-producer", p.Logger)
		if err != nil {
			panic(err)
		}
		return rabbitmq.NewConnector(consumer, producer, p.Router, p.Logger)
	case HTTP:
		panic("Need to be implemented")
	case FILE:
		panic("Need to be implemented")
	case CONSOLE:
		panic("Need to be implemented")
	default:
		panic("Unsupported Connector")
	}
}
