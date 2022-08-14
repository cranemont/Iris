package judge

import (
	"fmt"
	"time"
)

type compiler struct {
	sandbox Sandbox
}

func NewCompiler(sandbox Sandbox) *compiler {
	return &compiler{sandbox}
}

func (c *compiler) Compile(args CompileRequestDto) {
	fmt.Println("Compile! from Compiler")
	// go c.sandbox.Execute() // wait을 해야지, 아니면 sync로 돌리등가
	c.sandbox.Execute()
	time.Sleep(time.Second * 3)
	// 채널로 결과반환?
}
