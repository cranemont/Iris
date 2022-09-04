package egress

type Publisher interface {
	Publish(data string)
}
