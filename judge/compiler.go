package judge

import (
	"fmt"
	"strings"
	"time"

	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/judge/config"
)

type Compiler interface {
	Compile(out chan<- constants.GoResult, task *Task) // 얘는 task 몰라도 됨
}

type compiler struct {
	sandbox Sandbox
	config  *config.LanguageConfig
}

type CompileResult struct {
	Signal     int
	ErrorCode  int
	ExitCode   int
	ResultCode int
}

func NewCompiler(sandbox Sandbox, config *config.LanguageConfig) *compiler {
	return &compiler{sandbox, config}
}

func (c *compiler) Compile(out chan<- constants.GoResult, task *Task) {
	fmt.Println("Compile! from Compiler")

	options, err := c.config.Get(task.language) // 이게 된다고? private 아닌가? GetLanguage 가 필요없어?
	if err != nil {
		err := fmt.Errorf("failed to get language config: %s", err)
		out <- constants.GoResult{Err: err, Data: CompileResult{}}
		return
	}

	srcPath, err := c.config.GetSrcPath(task.dir, task.language)
	if err != nil {
		err := fmt.Errorf("failed to get language config: %s", err)
		out <- constants.GoResult{Err: err, Data: CompileResult{}}
		return
	}
	exePath, err := c.config.GetExePath(task.dir, task.language)
	if err != nil {
		err := fmt.Errorf("failed to get language config: %s", err)
		out <- constants.GoResult{Err: err, Data: CompileResult{}}
		return
	}

	// option에서 바로 매칭시켜서 sadnbox인자 넘겨주기

	args := strings.Replace(options.Args, "{srcPath}", srcPath, 1)
	args = strings.Replace(args, "{exePath}", exePath, 1)
	argSlice := strings.Split(args, " ")
	// sandbox 받지말고 그냥 여기서 arg처리한다음에 libjudger 실행하기

	c.sandbox.Execute(
		&SandboxArgs{
			ExePath:     options.CompilerPath,
			MaxCpuTime:  options.MaxCpuTime,
			MaxRealTime: options.MaxRealTime,
			MaxMemory:   options.MaxMemory,
			Args:        argSlice,
		})
	time.Sleep(time.Second * 2)
	// 채널로 결과반환?

	// sandbox result 추가
	// 컴파일 실패시 CompileResult에 error 추가
	out <- constants.GoResult{Err: err, Data: CompileResult{ResultCode: 0}}
}
