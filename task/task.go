package task

// Task가 여기있는게 맞나? Judger로 가야하는거 아닐까?
import (
	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/judger"
	"github.com/cranemont/judge-manager/mq"
	"github.com/cranemont/judge-manager/utils"
)

type Task struct {
	dir       string
	code      string
	language  string
	problemId string
	limits    mq.Limits
	testcases mq.Testcases
}

func NewTask(s mq.SubmissionDto) *Task {
	return &Task{
		dir:       utils.RandString(constants.DIR_NAME_LEN),
		code:      s.Code,
		language:  s.Language,
		problemId: s.ProblemId,
		limits:    s.Limits,
		testcases: s.Testcases,
	}
}

func (t *Task) GetDir() string {
	return t.dir
}

func (t *Task) ToCompileRequestDto() *judger.CompileRequestDto {
	return &judger.CompileRequestDto{Code: t.code, Language: t.language}
}

func (t *Task) ToJudgeRequestDto() *judger.JudgeRequestDto {
	runRequestDto := judger.RunRequestDto{}
	return &judger.JudgeRequestDto{RunRequestDto: &runRequestDto, Testcases: &t.testcases}
}
