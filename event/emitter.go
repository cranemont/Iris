package event

import "fmt"

type Emitter interface {
	Emit(eventName string, data interface{}) error
}

type emitter struct {
	eventMap map[string](chan interface{})
}

func NewEventEmitter(eventMap map[string](chan interface{})) *emitter {
	return &emitter{eventMap: eventMap}
}

func (e *emitter) Emit(eventName string, data interface{}) error {
	// https://stackoverflow.com/questions/2050391/how-to-check-if-a-map-contains-a-key-in-go
	if val, ok := e.eventMap[eventName]; ok {
		val <- data
		return nil
	}
	return fmt.Errorf("unregistered event: %s", eventName)
}
