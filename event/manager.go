package event

import (
	"github.com/cranemont/judge-manager/constants"
)

// handler 관리
// listener 등록
// 각각 객체생성

type EventManager struct {
	eventMap map[string](chan interface{})
	listener EventListener
	emitter  EventEmitter
}

func NewEventManager(
	eventMap map[string](chan interface{}),
	listener EventListener,
	emitter EventEmitter,
) *EventManager {
	return &EventManager{eventMap, listener, emitter}
}

// handlerFn?
func (e *EventManager) Listen(eventName string, handlerFn string) {
	ch := make(chan interface{}, constants.EVENT_CHAN_SIZE)
	e.eventMap[eventName] = ch

	go e.listener.On(ch, handlerFn)
}

func (e *EventManager) Dispatch(eventName string, data interface{}) {
	e.emitter.Emit(eventName, data)
}
