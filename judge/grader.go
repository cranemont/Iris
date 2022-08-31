package judge

import "fmt"

type Grader interface {
	Grade(answer []byte, output []byte) (bool, error)
}

type grader struct {
}

func NewGrader() *grader {
	return &grader{}
}

func (g *grader) Grade(answer []byte, output []byte) (bool, error) {
	// 일단 파일로 읽어서 채점
	// sed로 날리기
	// https://stackoverflow.com/questions/20521857/remove-white-space-from-the-end-of-line-in-linux

	fmt.Println("grading....")
	return true, nil
}
