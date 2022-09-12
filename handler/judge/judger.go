package judge

import (
	"errors"
	"fmt"

	"github.com/cranemont/judge-manager/common/exception"
	"github.com/cranemont/judge-manager/common/grade"
	"github.com/cranemont/judge-manager/common/result"
	"github.com/cranemont/judge-manager/sandbox"
	"github.com/cranemont/judge-manager/testcase"
)

var JudgeErrPrefix = "[judge: Judge]"
var ErrCompile = errors.New("compile failed")
var ErrTestcaseGet = errors.New("testcase get failed")

type Judger struct {
	compiler        sandbox.Compiler
	runner          sandbox.Runner
	testcaseManager testcase.Manager
}

func NewJudger(
	compiler sandbox.Compiler,
	runner sandbox.Runner,
	testcaseManager testcase.Manager,
) *Judger {
	return &Judger{
		compiler,
		runner,
		testcaseManager,
	}
}

func (j *Judger) Judge(task *JudgeTask) error {
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
		return fmt.Errorf("%s: %w", JudgeErrPrefix, compileOut.Err)
	}
	if testcaseOut.Err != nil {
		return fmt.Errorf("%s: %w: %s", JudgeErrPrefix, ErrTestcaseGet, testcaseOut.Err.Error())
	}

	compileResult, ok := compileOut.Data.(sandbox.CompileResult)
	if !ok {
		return fmt.Errorf("%s: %w: CompileResult", JudgeErrPrefix, exception.ErrTypeAssertionFail)
	}
	if compileResult.ResultCode != sandbox.SUCCESS {
		// 컴파일러를 실행했으나 컴파일에 실패한 경우
		task.CompileError(compileResult.ErrOutput)
		return fmt.Errorf("%s: %w", JudgeErrPrefix, ErrCompile)
	}

	tc, ok := testcaseOut.Data.(testcase.Testcase)
	if !ok {
		return fmt.Errorf("%s: %w: Testcase", JudgeErrPrefix, exception.ErrTypeAssertionFail)
	}

	tcNum := len(tc.Data)
	task.MakeRunResult(tcNum)

	// testcase의 order로 정렬
	// FIXME: out이라는 이름이 헷갈림 wrapper가 나을듯
	for i := 0; i < tcNum; i++ {
		res, err := j.runner.Run(sandbox.RunRequest{
			Order:       i,
			Dir:         task.dir,
			Language:    task.language,
			TimeLimit:   task.timeLimit,
			MemoryLimit: task.memoryLimit,
		}, []byte(tc.Data[i].In))
		if err != nil {
			task.SetRunResultCode(i, SYSTEM_ERROR)
			continue
		}
		task.SetRunResult(i, res)
		if res.ResultCode != sandbox.RUN_SUCCESS {
			continue
		}

		// 하나당 약 50microsec 10개 채점시 500microsec.
		// output이 커지면 더 길어짐 -> FIXME: 최적화 과정에서 goroutine으로 수정
		// st := time.Now()
		accepted, err := grade.Grade([]byte(tc.Data[i].Out), res.Output)
		if err != nil {
			task.SetRunResultCode(i, SYSTEM_ERROR)
			continue
		}
		if accepted {
			task.SetRunResultCode(i, ACCEPTED)
		} else {
			task.SetRunResultCode(i, WRONG_ANSWER)
		}
	}

	fmt.Println("done")
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
	fmt.Println(res)
	if err != nil {
		out <- result.ChResult{Err: err}
		return
	}
	out <- result.ChResult{Data: res}
}
