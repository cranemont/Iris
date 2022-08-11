package main

import (
	"github.com/cranemont/judger/judge"
	"github.com/cranemont/judger/judge/compiler"
	"github.com/cranemont/judger/judge/runner"
	"github.com/cranemont/judger/judge/sandbox"
)

func main() {

	sandbox := sandbox.NewSandbox()
	compiler := compiler.NewCompiler(sandbox)
	runner := runner.NewRunner(sandbox)
	processor := judge.NewProcessor(compiler, runner)
	processor.Judge()

	// 여기서 rabbitMQ consumer가 돌고
	// 메시지 수신시 채점자 호출
}
