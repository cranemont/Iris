package judger

import (
	"fmt"
	"time"

	"github.com/cranemont/judge-manager/judger/config"
)

type Compiler interface {
	Compile(dto *CompileRequestDto)
}

type compiler struct {
	sandbox Sandbox
	option  *config.CompileOption
}

func NewCompiler(sandbox Sandbox, option *config.CompileOption) *compiler {
	return &compiler{sandbox, option}
}

func (c *compiler) Compile(dto *CompileRequestDto) {
	fmt.Println("Compile! from Compiler")
	// go c.sandbox.Execute() // wait을 해야지, 아니면 sync로 돌리등가
	c.sandbox.Execute()
	time.Sleep(time.Second * 3)
	// 채널로 결과반환?
}
