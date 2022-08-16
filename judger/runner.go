package judger

import (
	"fmt"

	"github.com/cranemont/judge-manager/judger/config"
)

type Runner interface {
	Run(dto *RunRequestDto)
}

type runner struct {
	sandbox Sandbox
	option  *config.RunOption
}

func NewRunner(sandbox Sandbox, option *config.RunOption) *runner {
	return &runner{sandbox, option}
}

func (r *runner) Run(dto *RunRequestDto) {
	fmt.Println("RUN! from runner")
	r.sandbox.Execute()
	// 채널로 결과반환
}

func (r *runner) Result() {
	fmt.Println("Result is...!")
}
