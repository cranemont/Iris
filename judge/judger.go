package judge

import (
	"fmt"
	"strconv"

	"github.com/cranemont/judge-manager/common/dto"
	"github.com/cranemont/judge-manager/common/exception"
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
	grader          Grader
	testcaseManager testcase.TestcaseManager
}

func NewJudger(
	compiler sandbox.Compiler,
	runner sandbox.Runner,
	grader Grader,
	testcaseManager testcase.TestcaseManager,
) *Judger {
	return &Judger{
		compiler,
		runner,
		grader,
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

	// set testcase로 분리
	tc, ok := testcaseResult.Data.(testcase.Testcase)
	if !ok {
		return fmt.Errorf("%w: invalid testcase data", exception.ErrTypeAssertionFail)
	}
	tcNum := len(tc.Data)
	runOut := make(chan dto.GoResult, tcNum)
	for i := 0; i < tcNum; i++ {
		go j.run(runOut, task.dir, strconv.Itoa(i), task.language, nil)
	}

	gradeOut := make(chan dto.GoResult, tcNum)
	for i := 0; i < tcNum; i++ {
		result := <-runOut
		if t, ok := result.Data.(sandbox.RunResult); ok {
			fmt.Println(t.Output) // 이걸 아래 grade에 넘겨주기
		}
		go j.grade(gradeOut, nil, nil)
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

func (j *Judger) run(out chan<- dto.GoResult, dir string, id string, language string, input []byte) {
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
	result, err := j.grader.Grade(answer, output)
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
