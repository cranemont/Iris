package judge

import (
	"errors"
	"fmt"
	"log"

	"github.com/cranemont/judge-manager/common/exception"
	"github.com/cranemont/judge-manager/common/grade"
	"github.com/cranemont/judge-manager/common/result"
	"github.com/cranemont/judge-manager/sandbox"
	"github.com/cranemont/judge-manager/testcase"
)

var ErrCompileFail = errors.New("judge: compile failed")
var ErrRunFail = errors.New("judge: run failed")
var ErrGrade = errors.New("judge: grading failed")

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

	testcaseOutCh := make(chan result.ChResult)
	go j.getTestcase(testcaseOutCh, task.problemId)
	compileOutCh := make(chan result.ChResult)
	go j.compile(
		compileOutCh,
		sandbox.CompileRequest{
			Dir:      task.dir,
			Language: task.language,
		},
	)

	compileOut := <-compileOutCh
	testcaseOut := <-testcaseOutCh
	if compileOut.Err != nil {
		// NewError로 분리(funcName, error) 받아서 아래 포맷으로 에러 반환하는 함수
		// 컴파일러 실행 과정이나 이후 처리 과정에서 오류가 생긴 경우
		return fmt.Errorf("[judge: Judge]: %w", compileOut.Err)
	}
	if testcaseOut.Err != nil {
		return fmt.Errorf("[judge: Judge]: %w", testcaseOut.Err)
	}

	compileResult, ok := compileOut.Data.(sandbox.CompileResult)
	if !ok {
		return fmt.Errorf("%w: CompileResult", exception.ErrTypeAssertionFail)
	}
	if compileResult.ResultCode != sandbox.SUCCESS {
		// 컴파일러를 실행했으나 컴파일에 실패한 경우
		task.CompileError(compileResult.ErrOutput)
		return nil
	}

	tc, ok := testcaseOut.Data.(testcase.Testcase)
	if !ok {
		return fmt.Errorf("%w: Testcase", exception.ErrTypeAssertionFail)
	}

	// FIXME: 이 아래 과정 갈아엎기. Result를 중심으로
	tcNum := len(tc.Data)
	task.MakeRunResult(tcNum)

	runOutCh := make(chan result.ChResult, tcNum)
	for i := 0; i < tcNum; i++ {
		go j.run(
			runOutCh,
			sandbox.RunRequest{Order: i, Dir: task.dir, Language: task.language},
			[]byte(tc.Data[i].In),
		)
	}

	// FIXME: out이라는 이름이 헷갈림 wrapper가 나을듯
	gradeOutCh := make(chan result.ChResult, tcNum)
	for i := 0; i < tcNum; i++ {
		runOut := <-runOutCh
		order := runOut.Order
		if runOut.Err != nil {
			task.SetRunState(order, SYSTEM_ERROR)
			continue
			// RunData -> SYSTEM_ERROR
			// JudgeResult -> RUN_FAILED
		}
		runResult, ok := runOut.Data.(sandbox.RunResult)
		if !ok {
			task.SetRunState(order, SYSTEM_ERROR)
			log.Println("%w: RunResult", exception.ErrTypeAssertionFail)
			continue
		}
		task.SetRunResult(order, runResult)
		fmt.Print(order)
		go j.grade(gradeOutCh, []byte(tc.Data[order].Out), runResult.Output, order)
	}

	// FIXME: order 관리
	for i := 0; i < tcNum; i++ {
		gradeOut := <-gradeOutCh
		order := gradeOut.Order
		if gradeOut.Err != nil {
			task.SetRunState(order, SYSTEM_ERROR)
			continue
		}
		accepted, ok := gradeOut.Data.(bool)
		if !ok {
			task.SetRunState(order, SYSTEM_ERROR)
			log.Println("%w: GradeResult", exception.ErrTypeAssertionFail)
		}
		if accepted {
			task.SetRunState(order, ACCEPTED)
		} else {
			task.SetRunState(order, WRONG_ANSWER)
		}
	}
	// FIXME: 여기까지 수정

	fmt.Println("done")
	return nil
}

// wrapper to use goroutine
func (j *Judger) compile(out chan<- result.ChResult, dto sandbox.CompileRequest) {
	res, err := j.compiler.Compile(dto)
	if err != nil {
		out <- result.ChResult{Err: err}
	}
	out <- result.ChResult{Data: res}
}

func (j *Judger) run(out chan<- result.ChResult, dto sandbox.RunRequest, input []byte) {
	res, err := j.runner.Run(dto, nil)
	if err != nil {
		out <- result.ChResult{Err: err, Order: dto.Order}
	}
	out <- result.ChResult{Data: res, Order: dto.Order}
}

func (j *Judger) grade(out chan<- result.ChResult, answer []byte, output []byte, order int) {
	res, err := grade.Grade(answer, output)
	if err != nil {
		out <- result.ChResult{Err: err, Order: order}
	}
	out <- result.ChResult{Data: res, Order: order}
}

// wrapper to use goroutine
func (j *Judger) getTestcase(out chan<- result.ChResult, problemId string) {
	res, err := j.testcaseManager.GetTestcase(problemId)
	if err != nil {
		out <- result.ChResult{Err: err}
	}
	out <- result.ChResult{Data: res}
}
