package event

import (
	"github.com/cranemont/judge-manager/constants"
)

type Manager interface {
	Listen(eventName string, handlerFn string)
	Dispatch(eventName string, data interface{}) error
}

type manager struct {
	eventMap map[string](chan interface{})
	listener Listener
	emitter  Emitter
}

func NewManager(
	eventMap map[string](chan interface{}),
	listener Listener,
	emitter Emitter,
) *manager {
	return &manager{eventMap, listener, emitter}
}

// handlerFn?
func (e *manager) Listen(eventName string, handlerFn string) {
	ch := make(chan interface{}, constants.EVENT_CHAN_SIZE)
	e.eventMap[eventName] = ch

	go e.listener.On(ch, handlerFn)
}

func (e *manager) Dispatch(eventName string, data interface{}) error {
	err := e.emitter.Emit(eventName, data)
	if err != nil {
		return err
	}
	return nil
}
