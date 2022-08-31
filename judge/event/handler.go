package judgeEvent

import (
	"fmt"
	"log"

	"github.com/cranemont/judge-manager/common/exception"
	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/event"
	"github.com/cranemont/judge-manager/judge"
)

// controller의 역할. OnJudge, OnRun, OnOutput등으로 여러 상황 구분
type handler struct {
	funcMap      map[string](func(task *judge.Task) error)
	judgeService *judge.Judger
	eventEmitter event.Emitter
}

func NewHandler(
	judgeService *judge.Judger,
	eventEmitter event.Emitter,
) *handler {
	handler := &handler{judgeService: judgeService, eventEmitter: eventEmitter}
	funcMap := map[string](func(task *judge.Task) error){
		"OnExec": (*handler).OnExec,
		"OnExit": (*handler).OnExit,
	}
	handler.funcMap = funcMap
	return handler
}

// controller의 역할!
func (h *handler) OnExec(task *judge.Task) error {
	if err := h.judgeService.Judge(task); err != nil {
		return fmt.Errorf("onexec: %w", err)
	}
	// error 처리
	fmt.Println("triggerring event")
	if err := h.eventEmitter.Emit(constants.TASK_EXITED, task); err != nil {
		return fmt.Errorf("onexec: event emit failed: %w", err)
	}
	return nil
}

func (h *handler) OnExit(task *judge.Task) error {
	// 파일 삭제, task 결과 업데이트 등 정리작업
	if err := h.eventEmitter.Emit(constants.PUBLISH_RESULT, task); err != nil {
		return fmt.Errorf("onexit: event emit failed: %w", err)
	}
	return nil
}

func (h *handler) Call(funcName string, args interface{}) {
	//존재 확인. 없으면 registerFn 구현하라는 에러 throw
	if fn, ok := h.funcMap[funcName]; ok {
		if _, ok := args.(*judge.Task); ok {
			fmt.Println("handler function calling... ", funcName)
			if err := fn(args.(*judge.Task)); err != nil {
				log.Printf("error on %s: %s", funcName, err)
			}
		} else {
			log.Printf("%s: invalid task data", exception.ErrTypeAssertionFail)
		}
	} else {
		log.Printf("unregistered function: %s", funcName)
	}
}
