package judge

// Task가 NewTask로 생성되어야 하는 이유는?
import (
	"time"

	"github.com/cranemont/judge-manager/common/utils"
	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/ingress/rmq"
	"github.com/cranemont/judge-manager/sandbox"
)

// StatusCode
type RunData struct {
	Order      int    `json:"order"`
	ResultCode string `json:"resultCode"` // int for prod
	CpuTime    int    `json:"cpuTime"`
	RealTime   int    `json:"realTime"`
	Memory     int    `json:"memory"`
	Signal     int    `json:"signal"`
	ErrorCode  int    `json:"exitCode"`
	ExitCode   int    `json:"errorCode"`
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

type JudgeTaskResult struct {
	CompileErr string    `json:"compileError"`
	Run        []RunData `json:"runResult"`
}

// task interface, package? spj task, run task...
type JudgeTask struct {
	dir         string
	code        string
	language    string
	problemId   string
	timeLimit   int
	memoryLimit int
	Result      JudgeTaskResult
	StartedAt   time.Time // for time check
}

func NewTask(s rmq.JudgeRequest) *JudgeTask {
	// validate, initialize
	return &JudgeTask{
		dir:         utils.RandString(constants.DIR_NAME_LEN),
		code:        s.Code,
		language:    s.Language,
		problemId:   s.ProblemId,
		timeLimit:   s.TimeLimit,
		memoryLimit: s.MemoryLimit,
		Result:      JudgeTaskResult{},
	}
}

func (t *JudgeTask) GetDir() string {
	return t.dir
}

func (t *JudgeTask) GetCode() string {
	return t.code
}

func (t *JudgeTask) GetLanguage() string {
	return t.language
}

func (t *JudgeTask) CompileError(output string) {
	t.Result.CompileErr = output
}

func (t *JudgeTask) MakeRunResult(testcaseNum int) {
	t.Result.Run = make([]RunData, testcaseNum)
}

func (t *JudgeTask) SetRunResultCode(order int, stateCode string) {
	t.Result.Run[order].ResultCode = stateCode
}

func (t *JudgeTask) SetRunResult(order int, runResult sandbox.RunResult) {
	systemErr := false
	if runResult.ResultCode != sandbox.RUN_SUCCESS {
		switch runResult.ResultCode {
		case sandbox.CPU_TIME_LIMIT_EXCEEDED:
			t.SetRunResultCode(order, CPU_TLE)
		case sandbox.REAL_TIME_LIMIT_EXCEEDED:
			t.SetRunResultCode(order, REAL_TLE)
		case sandbox.MEMORY_LIMIT_EXCEEDED:
			t.SetRunResultCode(order, MEMORY_LIMIT_EXCEEDED)
		default:
			t.SetRunResultCode(order, SYSTEM_ERROR)
			systemErr = true
		}
	}
	if !systemErr {
		t.Result.Run[order].CpuTime = runResult.CpuTime
		t.Result.Run[order].RealTime = runResult.RealTime
		t.Result.Run[order].Memory = runResult.Memory
		t.Result.Run[order].Signal = runResult.Signal
		t.Result.Run[order].ErrorCode = runResult.ErrorCode
		t.Result.Run[order].ExitCode = runResult.ExitCode
	}
	// system error가 아니면 run result task에 반영(Resource usage)
}

// func (t *JudgeTask) ResultToJson() ([]byte, error) {
// 	data, err := json.Marshal(t.Result)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal JudgeResult: %w", err)
// 	}
// 	return data, nil
// }
