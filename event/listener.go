package event

import "fmt"

type EventListener interface {
	On(eventCh <-chan interface{}, handlerFn string)
}

type eventListener struct {
	eventMap map[string](chan interface{}) // 공유(emitter와)
	handler  EventHandler
}

func NewEventListener(eventMap map[string](chan interface{}), handler EventHandler) *eventListener {
	return &eventListener{eventMap, handler}
}

// type assertion is faster than reflection
// https://stackoverflow.com/questions/28024884/does-a-type-assertion-type-switch-have-bad-performance-is-slow-in-go
func (e *eventListener) On(eventCh <-chan interface{}, handlerFn string) {
	// TODO: handlerFn name으로 등록된 메서드 호출
	for {
		args := <-eventCh
		fmt.Println("Event Recv: ", handlerFn)
		// 여기서 goroutine으로 호출하기 때문에 나머지는 신경쓸 필요 없음. 모든 handler의 시작점은 하나의 goroutine
		go e.handler.Call(handlerFn, args)
	}
}
