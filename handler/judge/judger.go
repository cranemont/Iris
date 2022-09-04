package judge

import (
	"fmt"

	"github.com/cranemont/judge-manager/common/exception"
	"github.com/cranemont/judge-manager/common/result"
	"github.com/cranemont/judge-manager/handler/grade"
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
		return fmt.Errorf("%s: %w", errJudge, compileOut.Err)
	}
	if testcaseOut.Err != nil {
		task.SetStatus(TESTCASE_GET_FAILED)
		return fmt.Errorf("%s: %w", errJudge, testcaseOut.Err)
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
	runOutCh := make(chan result.ChResult, tcNum)
	for i := 0; i < tcNum; i++ {
		go j.run(
			runOutCh,
			sandbox.RunRequest{Order: i, Dir: task.dir, Language: task.language},
			[]byte(tc.Data[i].In),
		)
	}

	gradeOutCh := make(chan result.ChResult, tcNum) // out이라는 이름이 헷갈림 wrapper가 나을듯
	for i := 0; i < tcNum; i++ {
		runOut := <-runOutCh
		if runOut.Err != nil {
			task.SetStatus(RUN_FAILED)
			// FIXME: 이러면 최종결과의 해당 order는 비어있게됨. runOut? 에 order 정보 넣기
			return runOut.Err
			// RunData -> SYSTEM_ERROR
			// JudgeResult -> RUN_FAILED
		}
		runResult, ok := runOut.Data.(sandbox.RunResult)
		if !ok {
			return fmt.Errorf("%w: RunResult", exception.ErrTypeAssertionFail)
		}
		// run result task에 반영

		if runResult.ResultCode != sandbox.SUCCESS {
			// run result 저장
		}

		fmt.Print(runResult.Order)
		go j.grade(gradeOutCh, []byte(tc.Data[runResult.Order].Out), runResult.Output)
	}

	// FIXME: order 관리
	finalResult := []bool{}
	for i := 0; i < tcNum; i++ {
		gradeOut := <-gradeOutCh
		finalResult = append(finalResult, gradeOut.Data.(bool))
		// task에 결과 반영
	}
	// 여기까지

	fmt.Println(finalResult)
	task.SetStatus(SUCCESS) // RunData는 위에서 채움(run, grade하고 정보 채우기)

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
