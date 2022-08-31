package event

import (
	"log"
)

type Listener interface {
	On(eventCh <-chan interface{}, handlerFn string)
}

type listener struct {
	eventMap map[string](chan interface{}) // 공유(emitter와)
	handler  EventHandler
}

func NewListener(eventMap map[string](chan interface{}), handler EventHandler) *listener {
	return &listener{eventMap, handler}
}

// type assertion is faster than reflection
// https://stackoverflow.com/questions/28024884/does-a-type-assertion-type-switch-have-bad-performance-is-slow-in-go
func (e *listener) On(eventCh <-chan interface{}, handlerFn string) {
	for {
		args := <-eventCh
		log.Println("Event Recv: ", handlerFn) // event log
		// 여기서 goroutine으로 호출하기 때문에 나머지는 신경쓸 필요 없음. 모든 handler의 시작점은 하나의 goroutine
		// 따라서 handler는 return값이 없음
		go e.handler.Call(handlerFn, args)
	}
}
