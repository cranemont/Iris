package judge

import (
	"fmt"
)

type Judger interface {
	Compile(task *Task)
	Judge(task *Task) // run and grade
}

type judger struct {
	compiler Compiler
	runner   Runner
	grader   Grader
}

func NewJudger(compiler Compiler, runner Runner, grader Grader) *judger {
	return &judger{compiler, runner, grader}
}

// err 처리
func (j *judger) Compile(task *Task) {
	j.compiler.Compile(task)
}

// func (j *judger) Run(task *Task, out chan<- string)

// err 처리, Run이랑 Grade로 분리
func (j *judger) Judge(task *Task) {
	// run and grade
	tcNum := task.GetTestcase().GetTotal()

	ch := make(chan string, tcNum)
	for i := 0; i < tcNum; i++ {
		go j.runner.Run(task)
	}
	for i := 0; i < tcNum; i++ {
		result := <-ch
		fmt.Printf("%s Done!\n", result)
		// 여기서 이제 grade 고루틴으로 정리
	}
	close(ch)
	// close chan
}
