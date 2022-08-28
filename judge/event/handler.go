package judgeEvent

import (
	"fmt"

	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/event"
	"github.com/cranemont/judge-manager/judge"
)

// controller의 역할. OnJudge, OnRun, OnOutput등으로 여러 상황 구분
type handler struct {
	funcMap      map[string]func(task *judge.Task)
	judgeService *judge.JudgeService
	eventEmitter event.EventEmitter
}

func NewJudgeEventHandler(
	judgeService *judge.JudgeService,
	eventEmitter event.EventEmitter,
) *handler {
	funcMap := make(map[string]func(task *judge.Task), 2)

	// funcMap := map[string]func(h *handler, task *judge.Task){
	// 	"OnExec": (*handler).OnExec,
	// 	"OnExit": (*handler).OnExit,
	// }
	return &handler{funcMap, judgeService, eventEmitter}
}

func (h *handler) OnExec(task *judge.Task) {
	// 고루틴으로 JudgeHandler의 judge 호출
	// h.fileManager.CreateDir(task.GetDir())
	go h.judgeService.Judge(task)
}

func (h *handler) OnExit(task *judge.Task) {
	// 파일 삭제, task 결과 업데이트 등 정리작업
	h.eventEmitter.Emit(constants.PUBLISH_RESULT, task)
	// go h.fileManager.RemoveDir(task.GetDir())
}

func (h *handler) RegisterFn() {
	h.funcMap["OnExec"] = h.OnExec
	h.funcMap["OnExit"] = h.OnExit
}

func (h *handler) Call(funcName string, args interface{}) {
	//존재 확인. 없으면 registerFn 구현하라는 에러 throw
	if v, ok := args.(*judge.Task); ok {
		fmt.Println("handler function calling... ", funcName, v.GetDir())
	} else {
		// err log, return
		fmt.Println("error")
	}
	h.funcMap[funcName](args.(*judge.Task))
}