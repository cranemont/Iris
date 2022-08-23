package judge

import (
	"fmt"

	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/event"
	"github.com/cranemont/judge-manager/fileManager"
	"github.com/cranemont/judge-manager/testcase"
)

type JudgeService struct {
	judger          Judger
	eventEmitter    event.EventEmitter
	fileManager     *fileManager.FileManager
	testcaseManager testcase.TestcaseManager
}

// 여기서 파일생성, 삭제 관리
func NewJudgeService(
	judger Judger,
	eventEmitter event.EventEmitter,
	fileManager *fileManager.FileManager,
	testcaseManager testcase.TestcaseManager,
) *JudgeService {
	return &JudgeService{judger, eventEmitter, fileManager, testcaseManager}
}

func (j *JudgeService) Judge(task *Task) {
	defer func() {
		// 에러가 발생했다면? -> task error / error 이벤트 전송
		fmt.Println("triggerring event")
		j.eventEmitter.Emit(constants.TASK_EXITED, task)
	}()

	// 컴파일과 동시에 테스트케이스 가져오기(메모리에 올리기), 동시에 config에서 언어 설정 가져오기... 그것들을 task에 저장하기
	// task의 testcase가 있으면 isValid 체크한다음에 그거 쓰고, 없으면 가져와서 task의 testcase에 저장
	// 이후 m.judge 호출

	// go j.judger.Compile(task)
	j.judger.Compile(task)
	// go GetTestcase -> task에 있는지 먼저 검사
	// 오류나면 handler event listener에 넘기고 return
	// wait

	// err 처리
	j.judger.RunAndGrade(task)

	// eventManager한테 task done 이벤트 전송
	fmt.Println("done")
}
