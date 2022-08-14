package judge

import (
	"fmt"
)

type runner struct {
	sandbox Sandbox
}

func NewRunner(sandbox Sandbox) *runner {
	return &runner{sandbox}
}

func (r *runner) Run(args RunRequestDto) {
	fmt.Println("RUN! from runner")
	r.sandbox.Execute()
	// 채널로 결과반환
}

func (r *runner) Result() {
	fmt.Println("Result is...!")
}
