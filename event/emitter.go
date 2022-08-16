package event

type EventEmitter interface {
	Emit(eventName string, data interface{})
}

type eventEmitter struct {
	eventMap map[string](chan interface{})
}

func NewEventEmitter(eventMap map[string](chan interface{})) *eventEmitter {
	return &eventEmitter{eventMap: eventMap}
}

func (e *eventEmitter) Emit(eventName string, data interface{}) {
	// https://stackoverflow.com/questions/2050391/how-to-check-if-a-map-contains-a-key-in-go
	if val, ok := e.eventMap[eventName]; ok {
		val <- data
	}
	// err log. return
}
