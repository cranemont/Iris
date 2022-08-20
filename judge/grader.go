package judge

import "fmt"

type Grader interface {
	Grade(task *Task, out chan string)
}

type grader struct {
}

func NewGrader() *grader {
	return &grader{}
}

func (g *grader) Grade(task *Task, out chan string) {
	// 일단 파일로 읽어서 채점
	// sed로 날리기
	// https://stackoverflow.com/questions/20521857/remove-white-space-from-the-end-of-line-in-linux

	fmt.Println("grading....")
	out <- "task " + task.GetDir() + " done!\n"
}
