package egress

import "log"

type rmqPublisher struct {
}

func NewRmqPublisher() *rmqPublisher {
	return &rmqPublisher{}
}

func (r *rmqPublisher) Publish(data string) {
	log.Println("Task Publish!: ", data)
}
