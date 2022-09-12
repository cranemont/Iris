package sandbox

import (
	"encoding/json"
	"fmt"

	"github.com/cranemont/judge-manager/file"
)

type CompileResult struct {
	ResultCode int    // ?
	ErrOutput  string // compile error message
	ExecResult string // resource usage and metadata from sandbox
}

type CompileRequest struct {
	Dir      string
	Language string
}

type Compiler interface {
	Compile(dto CompileRequest) (CompileResult, error)
}

type compiler struct {
	langConfig LangConfig
	file       file.FileManager
}

func NewCompiler(langConfig LangConfig, file file.FileManager) *compiler {
	return &compiler{langConfig, file}
}

func (c *compiler) Compile(dto CompileRequest) (CompileResult, error) {
	fmt.Println("Compile! from Compiler")
	dir, language := dto.Dir, dto.Language
	fmt.Println(dir, language)

	execArgs, err := c.langConfig.ToCompileExecArgs(dir, language)
	if err != nil {
		return CompileResult{}, err
	}

	res, err := Exec(execArgs, nil)
	if err != nil {
		return CompileResult{}, err
	}

	compileResult := CompileResult{ResultCode: SUCCESS}
	if res.ResultCode != SUCCESS {
		sandboxResult, err := json.Marshal(res)
		if err != nil {
			return CompileResult{}, fmt.Errorf("invalid result format: %w", err)
		}

		compileOutputPath := c.file.MakeFilePath(dir, CompileOutFile).String()
		data, err := c.file.ReadFile(compileOutputPath)
		if err != nil {
			return CompileResult{}, fmt.Errorf("failed to read output file: %w", err)
		}
		// TODO: res.ErrorCode를 포함한 전체 output을 로그에 남기기
		compileResult.ResultCode = res.ResultCode
		compileResult.ExecResult = string(sandboxResult)
		compileResult.ErrOutput = string(data)
		fmt.Println("compile failed!: ", compileResult)
	}
	return compileResult, nil
}
