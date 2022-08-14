package judge

import (
	"fmt"
)

type sandbox struct {
}

func NewSandbox() *sandbox {
	return &sandbox{}
}

func (s *sandbox) Execute() {
	fmt.Println("Sandbox: func execute")
}
