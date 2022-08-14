package manager

import (
	"fmt"
	"time"

	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/judge"
	"github.com/cranemont/judge-manager/mq"
	"github.com/cranemont/judge-manager/utils"
)

type Task struct {
	dir         string
	code        string
	language    string
	problemId   string
	testcaseNum int
}

func NewTask(s mq.SubmissionDto) *Task {
	return &Task{
		dir:       utils.RandString(constants.DIR_NAME_LEN),
		code:      s.Code,
		language:  s.Language,
		problemId: s.ProblemId,
	}
}

func (t *Task) getDir() {
}

func (t *Task) GetTestcaseNum() int {
	return 3
}

func (t *Task) ToCompileRequestDto() judge.CompileRequestDto {
	return judge.CompileRequestDto{}
}

func (t *Task) ToRunRequestDto() judge.RunRequestDto {
	return judge.RunRequestDto{}
}

func (t *Task) RequestRun(input string, done chan string) {
	fmt.Printf("Running task %s. %s...\n", input, t.dir)
	time.Sleep(time.Second * 3)
	done <- t.problemId
}
