package judge

import (
	"errors"
	"fmt"

	"github.com/cranemont/judge-manager/common/exception"
	"github.com/cranemont/judge-manager/common/grade"
	"github.com/cranemont/judge-manager/common/result"
	"github.com/cranemont/judge-manager/logger"
	"github.com/cranemont/judge-manager/sandbox"
	"github.com/cranemont/judge-manager/testcase"
)

var ErrCompile = errors.New("compile failed")
var ErrTestcaseGet = errors.New("testcase get failed")

type Judger struct {
	compiler        sandbox.Compiler
	runner          sandbox.Runner
	testcaseManager testcase.Manager
	logging         *logger.Logger
}

func NewJudger(
	compiler sandbox.Compiler,
	runner sandbox.Runner,
	testcaseManager testcase.Manager,
	logging *logger.Logger,
) *Judger {
	return &Judger{
		compiler,
		runner,
		testcaseManager,
		logging,
	}
}

func (j *Judger) Judge(submissionId int, task *JudgeTask) error {
	j.logging.Debug("hander/judge: Judge - compile and get testcase")
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
		return fmt.Errorf("judge: %w", compileOut.Err)
	}
	if testcaseOut.Err != nil {
		return fmt.Errorf("judge: %w: %s", ErrTestcaseGet, testcaseOut.Err.Error())
	}

	compileResult, ok := compileOut.Data.(sandbox.CompileResult)
	if !ok {
		return fmt.Errorf("judge: %w: CompileResult", exception.ErrTypeAssertionFail)
	}
	if compileResult.ResultCode != sandbox.SUCCESS {
		// 컴파일러를 실행했으나 컴파일에 실패한 경우
		task.CompileError(compileResult.ErrOutput)
		return fmt.Errorf("judge: compile failed: %w", ErrCompile)
	}

	tc, ok := testcaseOut.Data.(testcase.Testcase)
	if !ok {
		return fmt.Errorf("judge: %w: Testcase", exception.ErrTypeAssertionFail)
	}

	tcNum := len(tc.Data)
	task.InitResult(submissionId, tcNum)

	j.logging.Debug("hander/judge: Judge - run and grade")
	for i := 0; i < tcNum; i++ {
		res, err := j.runner.Run(sandbox.RunRequest{
			Order:       i,
			Dir:         task.dir,
			Language:    task.language,
			TimeLimit:   task.timeLimit,
			MemoryLimit: task.memoryLimit,
		}, []byte(tc.Data[i].In))
		if err != nil {
			task.SetResultCode(i, SYSTEM_ERROR)
			continue
		}
		task.SetResult(i, tc.Data[i].Id, res)
		if res.ResultCode != sandbox.RUN_SUCCESS {
			continue
		}

		// 하나당 약 50microsec 10개 채점시 500microsec.
		// output이 커지면 더 길어짐 -> FIXME: 최적화 과정에서 goroutine으로 수정
		// st := time.Now()
		accepted := grade.Grade([]byte(tc.Data[i].Out), res.Output)
		if accepted {
			task.SetResultCode(i, ACCEPTED)
		} else {
			task.SetResultCode(i, WRONG_ANSWER)
		}

		// update RunResultCode on every iteration
		task.SetJudgeResultCode(i)
	}
	return nil
}

// wrapper to use goroutine
func (j *Judger) compile(out chan<- result.ChResult, dto sandbox.CompileRequest) {
	res, err := j.compiler.Compile(dto)
	if err != nil {
		out <- result.ChResult{Err: err}
		return
	}
	out <- result.ChResult{Data: res}
}

// wrapper to use goroutine
func (j *Judger) getTestcase(out chan<- result.ChResult, problemId string) {
	res, err := j.testcaseManager.GetTestcase(problemId)
	if err != nil {
		out <- result.ChResult{Err: err}
		return
	}
	out <- result.ChResult{Data: res}
}
