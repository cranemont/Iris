package runner

import (
	"fmt"

	"github.com/cranemont/judger/judge/sandbox"
)

type Runner struct {
	sandbox sandbox.Interface
}

func NewRunner(sandbox sandbox.Interface) *Runner {
	return &Runner{sandbox}
}

func (r Runner) Run(ch chan bool) {
	fmt.Println("RUN! from runner")
	go r.sandbox.Execute(ch)
}

func (r Runner) Result() {
	fmt.Println("Result is...!")
}
