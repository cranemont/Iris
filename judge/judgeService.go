package judge

import (
	"fmt"
	"log"

	"github.com/cranemont/judge-manager/fileManager"
	"github.com/cranemont/judge-manager/testcase"
)

type JudgeService struct {
	compiler        Compiler
	runner          Runner
	grader          Grader
	fileManager     fileManager.FileManager
	testcaseManager testcase.TestcaseManager
}

// 여기서 파일생성, 삭제 관리
func NewJudgeService(
	compiler Compiler,
	runner Runner,
	grader Grader,
	fileManager fileManager.FileManager,
	testcaseManager testcase.TestcaseManager,
) *JudgeService {
	return &JudgeService{
		compiler,
		runner,
		grader,
		fileManager,
		testcaseManager,
	}
}

func (j *JudgeService) Judge(task *Task) (err error) {
	defer func(err *error) {
		// 에러가 발생했다면? -> task error / error 이벤트 전송

		// go j.fileManager.RemoveDir(task.GetDir())
	}(&err)

	// 컴파일과 동시에 테스트케이스 가져오기(메모리에 올리기), 동시에 config에서 언어 설정 가져오기... 그것들을 task에 저장하기
	// task의 testcase가 있으면 isValid 체크한다음에 그거 쓰고, 없으면 가져와서 task의 testcase에 저장
	// 이후 m.judge 호출
	j.fileManager.CreateDir(task.GetDir())

	compileOut := make(chan int)
	testcaseOut := make(chan *testcase.Testcase)
	go j.CompileWithChannel(compileOut, task)
	go j.testcaseManager.GetTestcaseWithChannel(testcaseOut, task.problemId)

	compileResult := <-compileOut
	testcase := <-testcaseOut
	fmt.Println(testcase)
	task.testcase = *testcase

	if compileResult == -1 || testcase == nil {
		err = fmt.Errorf("TC or Compile Failed")
		return err
	}

	// err 처리
	j.RunAndGrade(task)

	// eventManager한테 task done 이벤트 전송
	fmt.Println("done")
	return nil
}

func (j *JudgeService) CompileWithChannel(out chan<- int, task *Task) {
	fmt.Println("COMPILE WITH CH")
	result, err := j.compiler.Compile(task)
	if err != nil {
		log.Println(err)
		out <- -1
	}
	out <- result
}

// err 처리, Run이랑 Grade로 분리
func (j *JudgeService) RunAndGrade(task *Task) {

	// run and grade
	tcNum := task.GetTestcase().Count()
	fmt.Println(tcNum)

	runCh := make(chan string, tcNum)
	for i := 0; i < tcNum; i++ {
		go j.runner.Run(task, runCh) // 여기서는 인자 정리해서 넘겨주기
	}

	gradeCh := make(chan string, tcNum)
	for i := 0; i < tcNum; i++ {
		result := <-runCh
		// result에 따라서 grade할지, 다른방식 쓸지 결정
		// run 결과를 파일로?
		fmt.Printf("%s grader running\n", result)
		go j.grader.Grade(task, gradeCh) // 여기서는 인자 정리해서 넘겨주기
		// 여기서 이제 grade 고루틴으로 정리
	}

	finalResult := ""
	for i := 0; i < tcNum; i++ {
		gradeResult := <-gradeCh
		finalResult += gradeResult
		// task에 결과 반영
	}

	fmt.Println(finalResult)
}
