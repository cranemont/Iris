package judge

import (
	"bytes"
	"fmt"
	"unicode"
)

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
	fmt.Printf("answer: %soutput: %s", string(answer), string(output))
	return bytes.Equal(answer, TrimWhitespaceBeforeNewline(output)), nil
}

func TrimWhitespaceBeforeNewline(a []byte) []byte {
	b := bytes.Split(a, []byte("\n"))

	for idx, val := range b {
		b[idx] = bytes.TrimRightFunc(val, unicode.IsSpace)
	}
	return bytes.Join(b, []byte("\n"))
}
