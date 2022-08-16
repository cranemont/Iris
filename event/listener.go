package event

import (
	"fmt"

	"github.com/cranemont/judge-manager/task"
)

type EventListener interface {
	On(eventCh <-chan interface{}, handlerCh string)
}

type taskEventListener struct {
	eventMap map[string](chan interface{}) // 공유(emitter와)

}

func NewTaskEventListener(eventMap map[string](chan interface{})) *taskEventListener {
	return &taskEventListener{eventMap: eventMap}
}

// you should make each eventListener for specific type
// this is because each listener has struct for the handlerFn
// this is also for the preformance because type assertion is faster than reflection
// https://stackoverflow.com/questions/28024884/does-a-type-assertion-type-switch-have-bad-performance-is-slow-in-go
func (e *taskEventListener) On(eventCh <-chan interface{}, handlerFn string) {
	// TODO: handlerFn name으로 등록된 메서드 호출
	for {
		data := <-eventCh
		v, ok := data.(*task.Task)
		if ok {
			fmt.Println(v.GetDir())
		} else {
			// err log, return
		}
		fmt.Println(handlerFn)
		fmt.Println("Event Recv")
	}
}
