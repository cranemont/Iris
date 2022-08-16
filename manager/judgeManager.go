package manager

import (
	"fmt"
	"sync"

	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/event"
	"github.com/cranemont/judge-manager/judger"
	"github.com/cranemont/judge-manager/task"
)

type JudgeManager struct {
	judger       judger.Judger
	eventEmitter event.EventEmitter
}

func NewJudgeManager(
	judger judger.Judger,
	eventEmitter event.EventEmitter,
	// fileManager
	// tcManager

	// task 결과 알림(정상 종료 혹은 error)을 전송할수있는 채널
	// error를 전송할수있는 채널
) *JudgeManager {
	return &JudgeManager{
		judger:       judger,
		eventEmitter: eventEmitter,
	}
}

func (j *JudgeManager) Exec(task *task.Task, wg *sync.WaitGroup) {
	defer func() {
		// 에러가 발생했다면? -> task error / error 이벤트 전송
		fmt.Println("clean up directory... trigger event")
		j.eventEmitter.Emit(constants.TASK_EXITED, task)
		wg.Done()
	}()

	// 디렉토리 만들고
	// 컴파일과 동시에 테스트케이스 가져오기(메모리에 올리기), 동시에 config에서 언어 설정 가져오기... 그것들을 task에 저장하기
	// task의 testcase가 있으면 isValid 체크한다음에 그거 쓰고, 없으면 가져와서 task의 testcase에 저장
	// 이후 m.judge 호출

	// task id로 dir 만들기
	// go j.judger.Compile(task.ToCompileRequestDto())
	j.judger.Compile(task.ToCompileRequestDto())
	// go GetTestcase -> task에 있는지 먼저 검사
	// 오류나면 handler event listener에 넘기고 return
	// wait

	// err 처리
	j.judger.Judge(task.ToJudgeRequestDto())

	// eventManager한테 task done 이벤트 전송
	fmt.Println("done")
}
