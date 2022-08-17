package judgeEvent

import (
	"fmt"

	"github.com/cranemont/judge-manager/judge"
)

type listener struct {
	eventMap map[string](chan interface{}) // 공유(emitter와)
	handler  *handler
}

func NewJudgeEventListener(eventMap map[string](chan interface{}), handler *handler) *listener {
	return &listener{eventMap, handler}
}

// you should make EventListener for each specific channel data types
// this is because eventListener has struct for the handlerFn
// this is also for the preformance because type assertion is faster than reflection
// https://stackoverflow.com/questions/28024884/does-a-type-assertion-type-switch-have-bad-performance-is-slow-in-go
func (e *listener) On(eventCh <-chan interface{}, handlerFn string) {
	// TODO: handlerFn name으로 등록된 메서드 호출
	for {
		data := <-eventCh
		v, ok := data.(*judge.Task)
		if ok {
			fmt.Println(v.GetDir())
		} else {
			// err log, return
		}

		// fmt.Println(handlerFn)
		// go handlerFn(data)
		fmt.Println("Event Recv")
	}
}
