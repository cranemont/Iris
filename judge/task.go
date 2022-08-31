package judge

// Task가 여기있는게 맞나? Judger로 가야하는거 아닐까?
import (
	"time"

	"github.com/cranemont/judge-manager/common/utils"
	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/mq"
)

type Task struct {
	dir       string
	code      string
	language  string
	problemId string
	limit     mq.Limit
	StartedAt time.Time // for time check
}

func NewTask(s mq.SubmissionDto) *Task {
	// validate, initialize
	return &Task{
		dir:       utils.RandString(constants.DIR_NAME_LEN),
		code:      s.Code,
		language:  s.Language,
		problemId: s.ProblemId,
		limit:     s.Limit,
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
