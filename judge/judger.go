package judge

import (
	"fmt"

	"github.com/cranemont/judge-manager/common/exception"
	"github.com/cranemont/judge-manager/common/result"
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

	testcaseOut := make(chan result.ChResult)
	go j.getTestcase(testcaseOut, task.problemId)
	compileOut := make(chan result.ChResult)
	go j.compile(
		compileOut,
		sandbox.CompileRequest{
			Dir:      task.dir,
			Language: task.language,
		},
	)

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
	runOut := make(chan result.ChResult, tcNum)
	for i := 0; i < tcNum; i++ {
		go j.run(
			runOut,
			sandbox.RunRequest{
				Id:       i,
				Dir:      task.dir,
				Language: task.language,
			},
			[]byte(tc.Data[i].In),
		)
	}

	gradeOut := make(chan result.ChResult, tcNum)
	for i := 0; i < tcNum; i++ {
		res := <-runOut
		runResult, ok := res.Data.(sandbox.RunResult)
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
func (j *Judger) compile(out chan<- result.ChResult, dto sandbox.CompileRequest) {
	// 여기서 결과값 처리
	res, err := j.compiler.Compile(dto)
	if err != nil {
		out <- result.ChResult{Err: fmt.Errorf("%s: %w", errCompile, err)}
	}
	// result 변환, 처리
	out <- result.ChResult{Data: res}
}

func (j *Judger) run(out chan<- result.ChResult, dto sandbox.RunRequest, input []byte) {
	// 여기서 결과값 처리
	res, err := j.runner.Run(dto, nil)
	if err != nil {
		out <- result.ChResult{Err: fmt.Errorf("%s: %w", errRun, err)}
	}
	// result 변환, 처리
	out <- result.ChResult{Data: res}
}

func (j *Judger) grade(out chan<- result.ChResult, answer []byte, output []byte) {
	// 여기서 결과값 처리
	res, err := grade.Grade(answer, output)
	if err != nil {
		out <- result.ChResult{Err: fmt.Errorf("%s: %w", errGrade, err)}
	}
	// result 변환, 처리
	out <- result.ChResult{Data: res}
}

// wrapper to use goroutine
func (j *Judger) getTestcase(out chan<- result.ChResult, problemId string) {
	// 여기서 결과값 처리
	res, err := j.testcaseManager.GetTestcase(problemId)
	if err != nil {
		out <- result.ChResult{Err: fmt.Errorf("%s: %w", errGetTestcase, err)}
	}
	// result 변환, 처리
	out <- result.ChResult{Data: res}
}
