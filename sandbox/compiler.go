package sandbox

import (
	"fmt"

	"github.com/cranemont/judge-manager/common/dto"
)

type Compiler interface {
	Compile(out chan<- dto.GoResult, dir string, language string) // 얘는 task 몰라도 됨
}

type compiler struct {
	sandbox Sandbox
	config  *LanguageConfig
}

type CompileResult struct {
	Signal     int
	ErrorCode  int
	ExitCode   int
	ResultCode int
}

func NewCompiler(sandbox Sandbox, config *LanguageConfig) *compiler {
	return &compiler{sandbox, config}
}

func (c *compiler) Compile(out chan<- dto.GoResult, dir string, language string) {
	fmt.Println("Compile! from Compiler")

	options, err := c.config.Get(language)
	if err != nil {
		out <- dto.GoResult{Err: err, Data: CompileResult{}}
		return
	}

	srcPath, err := c.config.MakeSrcPath(dir, language)
	if err != nil {
		out <- dto.GoResult{Err: err, Data: CompileResult{}}
		return
	}
	exePath, err := c.config.MakeExePath(dir, language)
	if err != nil {
		out <- dto.GoResult{Err: err, Data: CompileResult{}}
		return
	}
	argSlice, err := c.config.MakeArgSlice(srcPath, exePath, language)
	if err != nil {
		out <- dto.GoResult{Err: err, Data: CompileResult{}}
		return
	}

	c.sandbox.Execute(
		ExecArgs{
			ExePath:     options.CompilerPath,
			MaxCpuTime:  options.MaxCpuTime,
			MaxRealTime: options.MaxRealTime,
			MaxMemory:   options.MaxMemory,
			OutputPath:  "./compile/out.out",
			ErrorPath:   "./compile/error.out",
			LogPath:     "./compile/log.out",
			Args:        argSlice,
			Uid:         0,
			Gid:         0,
		}, nil)
	// time.Sleep(time.Second * 2)
	// 채널로 결과반환?

	// sandbox result 추가
	// 컴파일 실패시 CompileResult에 error 추가
	out <- dto.GoResult{Err: err, Data: CompileResult{ResultCode: 0}}
}