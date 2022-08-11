package judge

import (
	"github.com/cranemont/judger/judge/compiler"
	"github.com/cranemont/judger/judge/runner"
)

type Processor struct {
	compiler *compiler.Compiler
	runner   *runner.Runner
}

func NewProcessor(compiler *compiler.Compiler, runner *runner.Runner) *Processor {
	return &Processor{compiler: compiler, runner: runner}
}

func (p Processor) Judge() {
	ch := make(chan bool)
	p.compiler.Compile(ch)
	<-ch
	p.runner.Run(ch)
	<-ch
}
