package judge

import (
	"fmt"
	"log"
)

type Judger interface {
	Compile(task *Task)
	CompileWithChannel(task *Task, out chan<- int)
	RunAndGrade(task *Task) // run and grade
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
	j.compiler.Compile(task) // 여기서는 인자 정리해서 넘겨주기
	// TODO: 실패시 error 반환
}

func (j *judger) CompileWithChannel(task *Task, out chan<- int) {
	fmt.Println("COMPILE WITH CH")
	result, err := j.compiler.Compile(task)
	if err != nil {
		log.Println(err)
		out <- -1
	}
	out <- result
}

// func (j *judger) Run(task *Task, out chan<- string)

// err 처리, Run이랑 Grade로 분리
func (j *judger) RunAndGrade(task *Task) {

	// run and grade
	tcNum := task.GetTestcase().Count()
	fmt.Println(tcNum)

	runCh := make(chan string, tcNum)
	for i := 0; i < tcNum; i++ {
		go j.runner.Run(task, runCh) // 여기서는 인자 정리해서 넘겨주기
	}

	gradeCh := make(chan string, tcNum)
	for i := 0; i < tcNum; i++ {
		result := <-runCh
		// result에 따라서 grade할지, 다른방식 쓸지 결정
		// run 결과를 파일로?
		fmt.Printf("%s grader running\n", result)
		go j.grader.Grade(task, gradeCh) // 여기서는 인자 정리해서 넘겨주기
		// 여기서 이제 grade 고루틴으로 정리
	}

	finalResult := ""
	for i := 0; i < tcNum; i++ {
		gradeResult := <-gradeCh
		finalResult += gradeResult
		// task에 결과 반영
	}

	fmt.Println(finalResult)
}
