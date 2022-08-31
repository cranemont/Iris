package judge

import (
	"fmt"

	"github.com/cranemont/judge-manager/common/dto"
	"github.com/cranemont/judge-manager/common/exception"
	"github.com/cranemont/judge-manager/fileManager"
	"github.com/cranemont/judge-manager/sandbox"
	"github.com/cranemont/judge-manager/testcase"
)

type Judger struct {
	compiler        sandbox.Compiler
	runner          sandbox.Runner
	grader          Grader
	fileManager     fileManager.FileManager
	testcaseManager testcase.TestcaseManager
	config          *sandbox.LanguageConfig
}

// 여기서 파일생성, 삭제 관리
func NewJudger(
	compiler sandbox.Compiler,
	runner sandbox.Runner,
	grader Grader,
	fileManager fileManager.FileManager,
	testcaseManager testcase.TestcaseManager,
	config *sandbox.LanguageConfig,
) *Judger {
	return &Judger{
		compiler,
		runner,
		grader,
		fileManager,
		testcaseManager,
		config,
	}
}

func (j *Judger) Judge(task *Task) error {
	// 컴파일과 동시에 테스트케이스 가져오기(메모리에 올리기), 동시에 config에서 언어 설정 가져오기... 그것들을 task에 저장하기
	// task의 testcase가 있으면 isValid 체크한다음에 그거 쓰고, 없으면 가져와서 task의 testcase에 저장
	// 이후 m.judge 호출
	if err := j.fileManager.CreateDir(task.GetDir()); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	testcaseOut := make(chan dto.GoResult)
	go j.testcaseManager.GetTestcase(testcaseOut, task.problemId)

	srcPath, err := j.config.MakeSrcPath(task.dir, task.language)
	if err != nil {
		return err
	}
	if err := j.createSrcFile(srcPath, task.code); err != nil {
		return err
	}
	compileOut := make(chan dto.GoResult)
	go j.compiler.Compile(compileOut, task.dir, task.language)

	compileResult := <-compileOut
	testcaseResult := <-testcaseOut

	if compileResult.Err != nil {
		return fmt.Errorf("compile failed: %w", compileResult.Err)
	}
	if testcaseResult.Err != nil {
		return fmt.Errorf("testcase get failed: %w", testcaseResult.Err)
	}

	// set testcase로 분리
	if data, ok := testcaseResult.Data.(testcase.Testcase); ok {
		task.testcase = data
	} else {
		return fmt.Errorf("%w: invalid testcase data", exception.ErrTypeAssertionFail)
	}

	// err 처리
	j.RunAndGrade(task)

	// eventManager한테 task done 이벤트 전송
	fmt.Println("done")
	return nil
}

func (j *Judger) createSrcFile(srcPath string, code string) error {
	// task.code로 srcName에 파일 생성, 얘는 다른곳에서 생성해줘야됨. 컴파일이 아님
	if err := j.fileManager.CreateFile(srcPath, code); err != nil {
		// ENUM으로 변경, result code 반환
		err := fmt.Errorf("failed to create src file: %s", err)
		return err
	}
	return nil
}

// err 처리, Run이랑 Grade로 분리
func (j *Judger) RunAndGrade(task *Task) {

	// run and grade
	tcNum := task.GetTestcase().Count()
	fmt.Println(tcNum)

	runCh := make(chan dto.GoResult, tcNum)
	for i := 0; i < tcNum; i++ {
		go j.runner.Run(runCh, task.dir, task.language) // 여기서는 인자 정리해서 넘겨주기
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
