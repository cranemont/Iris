package sandbox

import (
	"fmt"
)

type Compiler interface {
	Compile(dir string, language string) (CompileResult, error) // 얘는 task 몰라도 됨
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

func (c *compiler) Compile(dir string, language string) (CompileResult, error) {
	fmt.Println("Compile! from Compiler")

	options, err := c.config.Get(language)
	if err != nil {
		return CompileResult{}, err
	}
	srcPath, err := c.config.MakeSrcPath(dir, language)
	if err != nil {
		return CompileResult{}, err
	}
	exePath, err := c.config.MakeExePath(dir, language)
	if err != nil {
		return CompileResult{}, err
	}
	argSlice, err := c.config.MakeArgSlice(srcPath, exePath, language)
	if err != nil {
		return CompileResult{}, err
	}

	result, err := c.sandbox.Execute(
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
		}, nil,
	)
	if err != nil {
		return CompileResult{}, err
	}
	fmt.Println(result)
	// time.Sleep(time.Second * 2)
	// 채널로 결과반환?

	// sandbox result 추가
	// 컴파일 실패시 CompileResult에 error 추가
	return CompileResult{ResultCode: 0}, nil
}
