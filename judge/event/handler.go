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
	funcMap      map[string]func(task *judge.Task)
	judgeService *judge.JudgeService
	eventEmitter event.Emitter
}

func NewJudgeEventHandler(
	judgeService *judge.JudgeService,
	eventEmitter event.Emitter,
) *handler {
	funcMap := make(map[string]func(task *judge.Task), 2)

	// funcMap := map[string]func(h *handler, task *judge.Task){
	// 	"OnExec": (*handler).OnExec,
	// 	"OnExit": (*handler).OnExit,
	// }
	return &handler{funcMap, judgeService, eventEmitter}
}

// controller의 역할!
func (h *handler) OnExec(task *judge.Task) {
	err := h.judgeService.Judge(task)
	if err != nil {
		log.Println("error onexec: %w", err)
		return
	}
	// error 처리
	fmt.Println("triggerring event")
	err = h.eventEmitter.Emit(constants.TASK_EXITED, task)
	if err != nil {
		log.Println("event emit failed: %w", err)
	}
}

func (h *handler) OnExit(task *judge.Task) {
	// 파일 삭제, task 결과 업데이트 등 정리작업
	err := h.eventEmitter.Emit(constants.PUBLISH_RESULT, task)
	if err != nil {
		log.Println("event emit failed: ", err)
	}
	// go h.fileManager.RemoveDir(task.GetDir())
}

func (h *handler) RegisterFn() {
	h.funcMap["OnExec"] = h.OnExec
	h.funcMap["OnExit"] = h.OnExit
}

func (h *handler) Call(funcName string, args interface{}) error {
	//존재 확인. 없으면 registerFn 구현하라는 에러 throw
	if _, ok := args.(*judge.Task); ok {
		fmt.Println("handler function calling... ", funcName)
		h.funcMap[funcName](args.(*judge.Task))
	} else {
		return fmt.Errorf("%w: invalid task data", exception.ErrTypeAssertionFail)
	}
	return nil
}
