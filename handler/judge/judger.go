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

var JudgeErrPrefix = "[judge: Judge]"
var ErrCompile = errors.New("compile failed")
var ErrTestcaseGet = errors.New("testcase get failed")

type Judger struct {
	testcaseManager testcase.TestcaseManager
}

func NewJudger(
	testcaseManager testcase.TestcaseManager,
) *Judger {
	return &Judger{
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

	// FIXME: 이 아래 과정 갈아엎기. Result, error를 중심으로
	tcNum := len(tc.Data)
	task.MakeRunResult(tcNum)

	runOutCh := make(chan result.ChResult, tcNum)
	// testcase의 order로 정렬
	for i := 0; i < tcNum; i++ {
		go j.run(
			runOutCh,
			sandbox.RunRequest{
				Order:       i,
				Dir:         task.dir,
				Language:    task.language,
				TimeLimit:   task.timeLimit,
				MemoryLimit: task.memoryLimit,
			},
			[]byte(tc.Data[i].In),
		)
	}

	// FIXME: out이라는 이름이 헷갈림 wrapper가 나을듯
	gradeNum := tcNum
	gradeOutCh := make(chan result.ChResult, tcNum)
	for i := 0; i < tcNum; i++ {
		runOut := <-runOutCh
		order := runOut.Order

		if runOut.Err != nil {
			task.SetRunResultCode(order, SYSTEM_ERROR)
			gradeNum -= 1
			continue
		}
		runResult, ok := runOut.Data.(sandbox.RunResult)

		if !ok {
			task.SetRunResultCode(order, SYSTEM_ERROR)
			log.Println("%w: RunResult", exception.ErrTypeAssertionFail)
			gradeNum -= 1
			continue
		}
		task.SetRunResult(order, runResult)
		fmt.Print(order)

		// result가 success가 아니면 grade 안함
		if runResult.ResultCode != sandbox.RUN_SUCCESS {
			gradeNum -= 1
			continue
		}
		go j.grade(gradeOutCh, []byte(tc.Data[order].Out), runResult.Output, order)
	}

	for i := 0; i < gradeNum; i++ {
		gradeOut := <-gradeOutCh
		order := gradeOut.Order

		if gradeOut.Err != nil {
			task.SetRunResultCode(order, SYSTEM_ERROR)
			continue
		}
		accepted, ok := gradeOut.Data.(bool)

		if !ok {
			task.SetRunResultCode(order, SYSTEM_ERROR)
			log.Println("%w: GradeResult", exception.ErrTypeAssertionFail)
		}
		if accepted {
			task.SetRunResultCode(order, ACCEPTED)
		} else {
			task.SetRunResultCode(order, WRONG_ANSWER)
		}
	}
	// FIXME: 여기까지 수정

	fmt.Println("done")
	return nil
}

// wrapper to use goroutine
func (j *Judger) compile(out chan<- result.ChResult, dto sandbox.CompileRequest) {
	res, err := sandbox.Compile(dto)
	if err != nil {
		out <- result.ChResult{Err: err}
		return
	}
	out <- result.ChResult{Data: res}
}

func (j *Judger) run(out chan<- result.ChResult, dto sandbox.RunRequest, input []byte) {
	res, err := sandbox.Run(dto, input)
	if err != nil {
		out <- result.ChResult{Err: err, Order: dto.Order}
		return
	}
	out <- result.ChResult{Data: res, Order: dto.Order}
}

func (j *Judger) grade(out chan<- result.ChResult, answer []byte, output []byte, order int) {
	res, err := grade.Grade(answer, output)
	if err != nil {
		out <- result.ChResult{Err: err, Order: order}
		return
	}
	out <- result.ChResult{Data: res, Order: order}
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
