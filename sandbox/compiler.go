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

func Compile(dto CompileRequest) (CompileResult, error) {
	fmt.Println("Compile! from Compiler")
	dir, language := dto.Dir, dto.Language
	fmt.Println(dir, language)

	languageConfig, err := GetConfig(language)
	if err != nil {
		return CompileResult{}, err
	}

	execArgs := languageConfig.ToCompileExecArgs(dir)
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
		data, err := file.ReadFile(languageConfig.CompileOutputPath(dir))
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
