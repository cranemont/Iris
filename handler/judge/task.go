package judge

// Task가 NewTask로 생성되어야 하는 이유는?
import (
	"time"

	"github.com/cranemont/judge-manager/common/utils"
	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/ingress/rmq"
)

// StatusCode
type RunData struct {
	Order      int
	StatusCode bool
	CpuTime    string
	RealTime   string
	Memory     string
	// Output? string
}

const ( // for RunData
	ACCEPTED = 0 + iota
	WRONG_ANSWER
	CPU_TLE
	REAL_TLE
	RUNTIME_ERROR
	SYSTEM_ERROR
)

type JudgeResult struct {
	// StatusCode int
	CompileErr string
	RunData    []RunData
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
		// 이걸 들고다니는게 맞을까?
		Result: JudgeResult{},
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

func (t *Task) ResultToJson() string {
	return "Judge Task to JSON"
}
