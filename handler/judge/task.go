package judge

// Task가 NewTask로 생성되어야 하는 이유는?
import (
	"time"

	"github.com/cranemont/judge-manager/common/utils"
	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/sandbox"
)

// const ( // for debug
// 	ACCEPTED                 = "accepted"
// 	WRONG_ANSWER             = "wrong answer"
// 	CPU_TIME_LIMIT_EXCEEDED  = "cpu time exceeded"
// 	REAL_TIME_LIMIT_EXCEEDED = "real time exceeded"
// 	MEMORY_LIMIT_EXCEEDED    = "memory exceeded"
// 	RUNTIME_ERROR            = "runtime error"
// 	SYSTEM_ERROR             = "system error"
// )

type JudgeTaskResult struct {
	SubmissionId    int         `json:"submissionId"`
	CompileErr      string      `json:"compileError"`
	TotalTestcase   int         `json:"totalTestcase"`
	AcceptedNum     int         `json:"acceptedNum"`
	JudgeResultCode int         `json:"judgeResultCode"` // first failed resultCode if some testcase failed
	Run             []RunResult `json:"runResult"`
}

type RunResult struct {
	TestcaseId string `json:"testcaseId"`
	ResultCode int    `json:"resultCode"` // int for prod
	CpuTime    int    `json:"cpuTime"`
	RealTime   int    `json:"realTime"`
	Memory     int    `json:"memory"`
	Signal     int    `json:"signal"`
	ErrorCode  int    `json:"exitCode"`
	ExitCode   int    `json:"errorCode"`
}

// RunResult ResultCode
const (
	ACCEPTED = 0 + iota
	WRONG_ANSWER
	CPU_TIME_LIMIT_EXCEEDED
	REAL_TIME_LIMIT_EXCEEDED
	MEMORY_LIMIT_EXCEEDED
	RUNTIME_ERROR
	SYSTEM_ERROR
)

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

func NewTask(
	code string,
	language string,
	problemId string,
	timeLimit int,
	memoryLimit int,
) *JudgeTask {
	// validate, initialize
	return &JudgeTask{
		dir:         utils.RandString(constants.DIR_NAME_LEN),
		code:        code,
		language:    language,
		problemId:   problemId,
		timeLimit:   timeLimit,
		memoryLimit: memoryLimit,
		Result:      JudgeTaskResult{},
		StartedAt:   time.Now(),
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

func (t *JudgeTask) InitResult(submissionId int, testcaseNum int) {
	t.Result.SubmissionId = submissionId
	t.Result.Run = make([]RunResult, testcaseNum)
	t.Result.TotalTestcase = testcaseNum
	t.Result.AcceptedNum = 0
}

func (t *JudgeTask) SetJudgeResultCode(idx int) {
	result := t.Result.Run[idx].ResultCode
	t.Result.JudgeResultCode = result
	if result == ACCEPTED {
		t.Result.AcceptedNum += 1
	}
}

func (t *JudgeTask) SetResultCode(idx int, stateCode int) {
	t.Result.Run[idx].ResultCode = stateCode
}

func (t *JudgeTask) SetResult(idx int, testcaseId string, runResult sandbox.RunResult) {
	systemErr := false
	if runResult.ResultCode != sandbox.RUN_SUCCESS {
		switch runResult.ResultCode {
		case sandbox.CPU_TIME_LIMIT_EXCEEDED:
			t.SetResultCode(idx, CPU_TIME_LIMIT_EXCEEDED)
		case sandbox.REAL_TIME_LIMIT_EXCEEDED:
			t.SetResultCode(idx, REAL_TIME_LIMIT_EXCEEDED)
		case sandbox.MEMORY_LIMIT_EXCEEDED:
			t.SetResultCode(idx, MEMORY_LIMIT_EXCEEDED)
		case sandbox.RUNTIME_ERROR:
			t.SetResultCode(idx, RUNTIME_ERROR)
		default:
			t.SetResultCode(idx, SYSTEM_ERROR)
			systemErr = true
		}
	}
	if !systemErr {
		t.Result.Run[idx].TestcaseId = testcaseId
		t.Result.Run[idx].CpuTime = runResult.CpuTime
		t.Result.Run[idx].RealTime = runResult.RealTime
		t.Result.Run[idx].Memory = runResult.Memory
		t.Result.Run[idx].Signal = runResult.Signal
		t.Result.Run[idx].ErrorCode = runResult.ErrorCode
		t.Result.Run[idx].ExitCode = runResult.ExitCode
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
