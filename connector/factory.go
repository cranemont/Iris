package connector

import (
	"os"

	"github.com/cranemont/judge-manager/connector/rabbitmq"
	"github.com/cranemont/judge-manager/router"
	"github.com/cranemont/judge-manager/service/logger"
)

type Connector interface {
	Connect()
	Disconnect()
	Handle(args ...any)
}

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
			os.Getenv("RABBITMQ_DEFAULT_USER") + ":" +
			os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" +
			os.Getenv("RABBITMQ_HOST") + ":" +
			os.Getenv("RABBITMQ_PORT") + "/"
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
