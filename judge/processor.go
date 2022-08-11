package judge

import (
	"github.com/cranemont/judger/judge/compiler"
	"github.com/cranemont/judger/judge/runner"
)

type Processor struct {
	compiler compiler.Interface
	runner   runner.Interface
}

func NewProcessor(compiler compiler.Interface, runner runner.Interface) *Processor {
	return &Processor{compiler, runner}
}

func (p Processor) Judge() {
	ch := make(chan bool)
	p.compiler.Compile(ch)
	<-ch
	p.runner.Run(ch)
	<-ch
}
