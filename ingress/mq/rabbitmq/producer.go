package rabbitmq

import (
	"log"
)

type Producer interface {
	Publish([]byte)
}

type producer struct {
}

func NewProducer() *producer {
	return &producer{}
}

func (p *producer) Publish(result []byte) {
	log.Printf("publisher: Publish: %s", string(result))
}
