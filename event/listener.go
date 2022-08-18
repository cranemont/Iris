package event

import "fmt"

type EventListener interface {
	On(eventCh <-chan interface{}, handlerFn string)
}

type listener struct {
	eventMap map[string](chan interface{}) // 공유(emitter와)
	handler  EventHandler
}

func NewEventListener(eventMap map[string](chan interface{}), handler EventHandler) *listener {
	return &listener{eventMap, handler}
}

// you should make EventListener for each specific channel data types
// this is because eventListener has struct for the handlerFn
// this is also for the preformance because type assertion is faster than reflection
// https://stackoverflow.com/questions/28024884/does-a-type-assertion-type-switch-have-bad-performance-is-slow-in-go
func (l *listener) On(eventCh <-chan interface{}, handlerFn string) {
	// TODO: handlerFn name으로 등록된 메서드 호출
	for {
		args := <-eventCh
		fmt.Println("Event Recv: ", handlerFn)
		go l.handler.Call(handlerFn, args)
	}
}
