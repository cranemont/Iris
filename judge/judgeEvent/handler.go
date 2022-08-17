package judgeEvent

import (
	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/event"
	"github.com/cranemont/judge-manager/fileManager"
	"github.com/cranemont/judge-manager/judge"
)

type handler struct {
	// taskMap      map[string](*task.Task)
	judgeController *judge.JudgeController
	fileManager     *fileManager.FileManager
	eventEmitter    event.EventEmitter
}

func NewJudgeEventHandler(
	judgeController *judge.JudgeController,
	fileManager *fileManager.FileManager,
	eventEmitter event.EventEmitter,
) *handler {
	return &handler{judgeController, fileManager, eventEmitter}
}

func (t *handler) OnExec(task *judge.Task) {
	// 고루틴으로 JudgeHandler의 judge 호출
	t.fileManager.CreateDir(task.GetDir())
	go t.judgeController.Judge(task)
}

func (t *handler) OnExit(task *judge.Task) {
	// 파일 삭제, task 결과 업데이트 등 정리작업
	t.eventEmitter.Emit(constants.PUBLISH_RESULT, task)
	go t.fileManager.RemoveDir(task.GetDir())
}
