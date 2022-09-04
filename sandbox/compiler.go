package sandbox

import (
	"encoding/json"
	"fmt"

	"github.com/cranemont/judge-manager/file"
)

type Compiler interface {
	Compile(dto CompileRequest) (CompileResult, error) // 얘는 task 몰라도 됨
}

type compiler struct {
	config *LanguageConfig
}

type CompileResult struct {
	ResultCode int    // ?
	ErrOutput  string // compile error message
	ExecResult string // resource usage and metadata from sandbox
}

type CompileRequest struct {
	Dir      string
	Language string
}

func NewCompiler(config *LanguageConfig) *compiler {
	return &compiler{config}
}

func (c *compiler) Compile(dto CompileRequest) (CompileResult, error) {
	fmt.Println("Compile! from Compiler")
	dir, language := dto.Dir, dto.Language
	fmt.Println(dir, language)

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
	argSlice, err := c.config.MakeCompileArgSlice(srcPath, exePath, language)
	if err != nil {
		return CompileResult{}, err
	}

	outputPath := file.MakeFilePath(dir, "compile.out").String()
	res, err := Exec(
		ExecArgs{
			ExePath:       options.CompilerPath,
			MaxCpuTime:    options.MaxCpuTime,
			MaxRealTime:   options.MaxRealTime,
			MaxMemory:     options.MaxMemory,
			MaxStackSize:  128 * 1024 * 1024,
			MaxOutputSize: 20 * 1024 * 1024,
			OutputPath:    outputPath,
			ErrorPath:     outputPath,
			LogPath:       "./log/compile/log.out",
			Args:          argSlice,
		}, nil,
	)
	// Exec fail
	if err != nil {
		return CompileResult{}, err
	}

	compileResult := CompileResult{ResultCode: SUCCESS}
	if res.ResultCode != SUCCESS {
		sandboxResult, err := json.Marshal(res)
		if err != nil {
			return CompileResult{}, fmt.Errorf("invalid result format: %w", err)
		}
		data, err := file.ReadFile(outputPath)
		if err != nil {
			return CompileResult{}, fmt.Errorf("failed to read output file: %w", err)
		}
		// TODO: res.ErrorCode를 포함한 전체 output을 로그에 남기기
		compileResult.ResultCode = res.ResultCode
		compileResult.ExecResult = string(sandboxResult)
		compileResult.ErrOutput = string(data)
		fmt.Println(compileResult)
	}
	// time.Sleep(time.Second * 2)
	// 채널로 결과반환?

	// sandbox result 추가
	// 컴파일 실패시 CompileResult에 error 추가
	return compileResult, nil
}
