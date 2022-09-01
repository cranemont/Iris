package judge

// Task가 NewTask로 생성되어야 하는 이유는?
import (
	"time"

	"github.com/cranemont/judge-manager/common/utils"
	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/mq"
)

type Task struct {
	dir         string
	code        string
	language    string
	problemId   string
	timeLimit   int
	memoryLimit int
	Status      string
	StartedAt   time.Time // for time check
}

func NewTask(s mq.SubmissionDto) *Task {
	// validate, initialize
	return &Task{
		dir:         utils.RandString(constants.DIR_NAME_LEN),
		code:        s.Code,
		language:    s.Language,
		problemId:   s.ProblemId,
		timeLimit:   s.TimeLimit,
		memoryLimit: s.MemoryLimit,
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
