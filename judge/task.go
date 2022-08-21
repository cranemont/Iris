package judge

// Task가 여기있는게 맞나? Judger로 가야하는거 아닐까?
import (
	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/mq"
	tc "github.com/cranemont/judge-manager/testcase"
	"github.com/cranemont/judge-manager/utils"
)

type Task struct {
	dir       string
	code      string
	language  string
	problemId string
	limit     mq.Limit
	testcase  tc.Testcase
}

func NewTask(s mq.SubmissionDto) *Task {
	// validate, initialize
	return &Task{
		dir:       utils.RandString(constants.DIR_NAME_LEN),
		code:      s.Code,
		language:  s.Language,
		problemId: s.ProblemId,
		limit:     s.Limit,
		testcase:  s.Testcase,
	}
}

func (t *Task) GetDir() string {
	return t.dir
}

func (t *Task) GetTestcase() *tc.Testcase {
	return &t.testcase
}
