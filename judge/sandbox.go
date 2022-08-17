package judge

import (
	"fmt"
)

type Sandbox interface {
	Execute()
}

type sandbox struct {
}

func NewSandbox() *sandbox {
	return &sandbox{}
}

func (s *sandbox) Execute() {
	fmt.Println("Sandbox: func execute")
	// stdin, out 연결해서 실행?
}
