package grade

import (
	"bytes"
	"fmt"
	"unicode"
)

// type Grader interface {
// 	Grade(answer []byte, output []byte) (bool, error)
// }

// type grader struct {
// }

// func NewGrader() *grader {
// 	return &grader{}
// }

func Grade(answer []byte, output []byte) (bool, error) {
	fmt.Println("grading....")
	// fmt.Printf("answer: %soutput: %s", string(answer), string(output))
	return bytes.Equal(answer, TrimWhitespaceBeforeNewline(output)), nil
}

func TrimWhitespaceBeforeNewline(a []byte) []byte {
	b := bytes.Split(a, []byte("\n"))

	for idx, val := range b {
		b[idx] = bytes.TrimRightFunc(val, unicode.IsSpace)
	}
	return bytes.Join(b, []byte("\n"))
}
