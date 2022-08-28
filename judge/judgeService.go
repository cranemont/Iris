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
	fileManager     fileManager.FileManager
	testcaseManager testcase.TestcaseManager
}

// 여기서 파일생성, 삭제 관리
func NewJudgeService(
	judger Judger,
	eventEmitter event.EventEmitter,
	fileManager fileManager.FileManager,
	testcaseManager testcase.TestcaseManager,
) *JudgeService {
	return &JudgeService{judger, eventEmitter, fileManager, testcaseManager}
}

func (j *JudgeService) Judge(task *Task) (err error) {
	defer func(err *error) {
		// 에러가 발생했다면? -> task error / error 이벤트 전송
		if *err != nil {
			fmt.Println("Error on judgeService.Judge: ", *err)
		} else {
			fmt.Println("triggerring event")
			j.eventEmitter.Emit(constants.TASK_EXITED, task)
		}
		// go j.fileManager.RemoveDir(task.GetDir())
	}(&err)

	// 컴파일과 동시에 테스트케이스 가져오기(메모리에 올리기), 동시에 config에서 언어 설정 가져오기... 그것들을 task에 저장하기
	// task의 testcase가 있으면 isValid 체크한다음에 그거 쓰고, 없으면 가져와서 task의 testcase에 저장
	// 이후 m.judge 호출
	j.fileManager.CreateDir(task.GetDir())

	compileCh := make(chan int)
	testcaseCh := make(chan *testcase.Testcase)
	go j.judger.CompileWithChannel(task, compileCh)
	go j.testcaseManager.GetTestcaseWithChannel(task.problemId, testcaseCh)

	compileResult := <-compileCh
	testcase := <-testcaseCh
	fmt.Println(testcase)
	task.testcase = *testcase

	if compileResult == -1 || testcase == nil {
		err = fmt.Errorf("TC or Compile Failed")
		return err
	}

	// err 처리
	j.judger.RunAndGrade(task)

	// eventManager한테 task done 이벤트 전송
	fmt.Println("done")
	return nil
}
