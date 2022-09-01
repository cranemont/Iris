package judge

import (
	"fmt"

	"github.com/cranemont/judge-manager/common/dto"
	"github.com/cranemont/judge-manager/common/exception"
	"github.com/cranemont/judge-manager/judge/grade"
	"github.com/cranemont/judge-manager/sandbox"
	"github.com/cranemont/judge-manager/testcase"
)

var errJudge = "[Judger: Judge]"
var errCompile = "[Judger: compile]"
var errRun = "[Judger: run]"
var errGrade = "[Judger: Grade]"
var errGetTestcase = "[Judger: getTestcase]"

type Judger struct {
	compiler        sandbox.Compiler
	runner          sandbox.Runner
	testcaseManager testcase.TestcaseManager
}

func NewJudger(
	compiler sandbox.Compiler,
	runner sandbox.Runner,
	testcaseManager testcase.TestcaseManager,
) *Judger {
	return &Judger{
		compiler,
		runner,
		testcaseManager,
	}
}

func (j *Judger) Judge(task *Task) error {
	// testcase 있는건 다른 함수에서 처리. grade가 필요없는 요청임

	testcaseOut := make(chan dto.GoResult)
	go j.getTestcase(testcaseOut, task.problemId)
	compileOut := make(chan dto.GoResult)
	go j.compile(compileOut, task.dir, task.language)

	compileResult := <-compileOut
	testcaseResult := <-testcaseOut
	if compileResult.Err != nil {
		// NewError로 분리(funcName, error) 받아서 아래 포맷으로 에러 반환하는 함수
		return fmt.Errorf("%s: %w", errJudge, compileResult.Err)
	}
	if testcaseResult.Err != nil {
		return fmt.Errorf("%s: %w", errJudge, testcaseResult.Err)
	}

	tc, ok := testcaseResult.Data.(testcase.Testcase)
	if !ok {
		return fmt.Errorf("%w: invalid testcase data", exception.ErrTypeAssertionFail)
	}
	tcNum := len(tc.Data)

	// 이 아래 과정 너무 지저분함. result 확인 과정은 wrapper function으로 넘길것
	runOut := make(chan dto.GoResult, tcNum)
	for i := 0; i < tcNum; i++ {
		go j.run(runOut, task.dir, i, task.language, []byte(tc.Data[i].In))
	}

	gradeOut := make(chan dto.GoResult, tcNum)
	for i := 0; i < tcNum; i++ {
		result := <-runOut
		runResult, ok := result.Data.(sandbox.RunResult)
		if !ok {
			return fmt.Errorf("%w: invalid RunResult data", exception.ErrTypeAssertionFail)
		}
		// runResult가 정상이라면
		fmt.Print(runResult.Id)
		go j.grade(gradeOut, []byte(tc.Data[runResult.Id].Out), runResult.Output)
	}

	finalResult := []bool{}
	for i := 0; i < tcNum; i++ {
		gradeResult := <-gradeOut
		finalResult = append(finalResult, gradeResult.Data.(bool))
		// task에 결과 반영
	}

	fmt.Println(finalResult)

	// eventManager한테 task done 이벤트 전송
	fmt.Println("done")
	return nil
}

// wrapper to use goroutine
func (j *Judger) compile(out chan<- dto.GoResult, dir string, language string) {
	// 여기서 결과값 처리
	result, err := j.compiler.Compile(dir, language)
	if err != nil {
		out <- dto.GoResult{Err: fmt.Errorf("%s: %w", errCompile, err)}
	}
	// result 변환, 처리
	out <- dto.GoResult{Data: result}
}

func (j *Judger) run(out chan<- dto.GoResult, dir string, id int, language string, input []byte) {
	// 여기서 결과값 처리
	result, err := j.runner.Run(dir, id, language, nil)
	if err != nil {
		out <- dto.GoResult{Err: fmt.Errorf("%s: %w", errRun, err)}
	}
	// result 변환, 처리
	out <- dto.GoResult{Data: result}
}

func (j *Judger) grade(out chan<- dto.GoResult, answer []byte, output []byte) {
	// 여기서 결과값 처리
	result, err := grade.Grade(answer, output)
	if err != nil {
		out <- dto.GoResult{Err: fmt.Errorf("%s: %w", errGrade, err)}
	}
	// result 변환, 처리
	out <- dto.GoResult{Data: result}
}

// wrapper to use goroutine
func (j *Judger) getTestcase(out chan<- dto.GoResult, problemId string) {
	// 여기서 결과값 처리
	result, err := j.testcaseManager.GetTestcase(problemId)
	if err != nil {
		out <- dto.GoResult{Err: fmt.Errorf("%s: %w", errGetTestcase, err)}
	}
	// result 변환, 처리
	out <- dto.GoResult{Data: result}
}
