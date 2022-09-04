package judge

// Task가 NewTask로 생성되어야 하는 이유는?
import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cranemont/judge-manager/common/utils"
	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/ingress/rmq"
	"github.com/cranemont/judge-manager/sandbox"
)

// StatusCode
type RunData struct {
	Order     int    `json:"order"`
	StateCode string `json:"stateCode"` // int for prod
	CpuTime   int    `json:"cpuTime"`
	RealTime  int    `json:"realTime"`
	Memory    int    `json:"memory"`
}

// RunData Status Code
// const (
// 	ACCEPTED = 0 + iota
// 	WRONG_ANSWER
// 	CPU_TLE
// 	REAL_TLE
// 	MEMORY_LIMIT_EXCEEDED
// 	RUNTIME_ERROR
// 	SYSTEM_ERROR
// )

const ( // for debug
	ACCEPTED              = "accepted"
	WRONG_ANSWER          = "wrong answer"
	CPU_TLE               = "cpu time exceeded"
	REAL_TLE              = "real time exceeded"
	MEMORY_LIMIT_EXCEEDED = "memory exceeded"
	RUNTIME_ERROR         = "runtime error"
	SYSTEM_ERROR          = "system error"
)

type JudgeResult struct {
	// StatusCode int
	CompileErr string    `json:"compileError"`
	Run        []RunData `json:"runResult"`
}

// task interface, package? spj task, run task...
type Task struct {
	dir         string
	code        string
	language    string
	problemId   string
	timeLimit   int
	memoryLimit int
	Result      JudgeResult
	StartedAt   time.Time // for time check
}

func NewTask(s rmq.JudgeRequest) *Task {
	// validate, initialize
	return &Task{
		dir:         utils.RandString(constants.DIR_NAME_LEN),
		code:        s.Code,
		language:    s.Language,
		problemId:   s.ProblemId,
		timeLimit:   s.TimeLimit,
		memoryLimit: s.MemoryLimit,
		Result:      JudgeResult{},
	}
}

func (t *Task) GetDir() string {
	return t.dir
}

func (t *Task) GetCode() string {
	return t.code
}

func (t *Task) GetLanguage() string {
	return t.language
}

func (t *Task) CompileError(output string) {
	t.Result.CompileErr = output
}

func (t *Task) MakeRunResult(testcaseNum int) {
	t.Result.Run = make([]RunData, testcaseNum)
}

func (t *Task) SetRunState(order int, stateCode string) {
	t.Result.Run[order].StateCode = stateCode
}

func (t *Task) SetRunResult(order int, runResult sandbox.RunResult) {
	systemErr := false
	if runResult.ResultCode != sandbox.RUN_SUCCESS {
		switch runResult.ResultCode {
		case sandbox.CPU_TIME_LIMIT_EXCEEDED:
			t.SetRunState(order, CPU_TLE)
		case sandbox.REAL_TIME_LIMIT_EXCEEDED:
			t.SetRunState(order, REAL_TLE)
		case sandbox.MEMORY_LIMIT_EXCEEDED:
			t.SetRunState(order, MEMORY_LIMIT_EXCEEDED)
		default:
			t.SetRunState(order, SYSTEM_ERROR)
			systemErr = true
		}
	}
	if !systemErr {
		t.Result.Run[order].CpuTime = runResult.CpuTime
		t.Result.Run[order].RealTime = runResult.RealTime
		t.Result.Run[order].Memory = runResult.Memory
	}
	// system error가 아니면 run result task에 반영(Resource usage)
}

func (t *Task) ResultToJson() (string, error) {
	data, err := json.Marshal(t.Result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JudgeResult: %w", err)
	}
	return string(data), nil
}
