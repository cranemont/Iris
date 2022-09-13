package sandbox

import (
	"fmt"
	"strconv"

	"github.com/cranemont/judge-manager/file"
)

type RunResult struct {
	Order      int
	ResultCode int
	ErrOutput  string // []byte?
	CpuTime    int
	RealTime   int
	Memory     int
	Signal     int
	ErrorCode  int
	ExitCode   int
	Output     []byte
}

type RunRequest struct {
	Order       int
	Dir         string
	Language    string
	TimeLimit   int
	MemoryLimit int
}

type Runner interface {
	Run(dto RunRequest, input []byte) (RunResult, error)
}

type runner struct {
	sandbox    Sandbox
	langConfig LangConfig
	file       file.FileManager
}

func NewRunner(sandbox Sandbox, langConfig LangConfig, file file.FileManager) *runner {
	return &runner{sandbox, langConfig, file}
}

func (r *runner) Run(dto RunRequest, input []byte) (RunResult, error) {
	fmt.Println("RUN! from runner")
	dir := dto.Dir
	order := dto.Order
	language := dto.Language
	timeLimit := dto.TimeLimit
	memoryLimit := dto.MemoryLimit

	// languageConfig, err := r.config.GetConfig(language)
	// if err != nil {
	// 	return RunResult{}, err
	// }

	execArgs, err := r.langConfig.ToRunExecArgs(
		dir,
		language,
		order,
		Limit{
			CpuTime:  timeLimit,
			RealTime: timeLimit * 3,
			Memory:   memoryLimit,
		},
		false,
	)
	if err != nil {
		return RunResult{}, err
	}

	res, err := r.sandbox.Exec(execArgs, input)
	if err != nil {
		return RunResult{}, fmt.Errorf("runner: Run failed: %w", err)
	}

	runResult := RunResult{
		Order:      order,
		ResultCode: SUCCESS,
		CpuTime:    res.CpuTime,
		RealTime:   res.RealTime,
		Memory:     res.Memory,
		Signal:     res.Signal,
		ErrorCode:  res.ErrorCode,
		ExitCode:   res.ExitCode,
	}

	outputPath := r.file.MakeFilePath(dir, strconv.Itoa(order)+".out").String()
	errorPath := r.file.MakeFilePath(dir, strconv.Itoa(order)+".error").String()

	if res.ResultCode != SUCCESS {
		// TODO: 함수로 분리
		data, err := r.file.ReadFile(errorPath)
		if err != nil {
			return RunResult{}, fmt.Errorf("runner: failed to read error file: %w", err)
		}
		runResult.ResultCode = res.ResultCode
		runResult.ErrOutput = string(data)
		fmt.Println(res) // log?
	}

	data, err := r.file.ReadFile(outputPath)
	if err != nil {
		return RunResult{}, fmt.Errorf("runner: failed to read output file: %w", err)
	}
	runResult.Output = data

	// fmt.Println(res)
	return runResult, nil
}
