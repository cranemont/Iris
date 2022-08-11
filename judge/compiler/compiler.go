package compiler

import (
	"fmt"

	"github.com/cranemont/judger/judge/sandbox"
)

type Compiler struct {
	sandbox sandbox.Interface
}

func NewCompiler(sandbox sandbox.Interface) *Compiler {
	return &Compiler{sandbox}
}

func (c *Compiler) Compile(ch chan bool) {
	fmt.Println("Compile! from Compiler")
	go c.sandbox.Execute(ch)
}
