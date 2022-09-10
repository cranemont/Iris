package sandbox

import (
	"fmt"

	"github.com/cranemont/judge-manager/common/file"
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

func Run(dto RunRequest, input []byte) (RunResult, error) {
	fmt.Println("RUN! from runner")
	dir := dto.Dir
	order := dto.Order
	language := dto.Language
	timeLimit := dto.TimeLimit
	memoryLimit := dto.MemoryLimit

	languageConfig, err := GetConfig(language)
	if err != nil {
		return RunResult{}, err
	}

	execArgs := languageConfig.ToRunExecArgs(
		dir,
		order,
		Limit{
			CpuTime:  timeLimit,
			RealTime: timeLimit * 3,
			Memory:   memoryLimit,
		},
		false,
	)

	res, err := Exec(execArgs, input)
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

	outputPath := languageConfig.RunOutputPath(dir, order)
	errorPath := languageConfig.RunErrPath(dir, order)

	if res.ResultCode != SUCCESS {
		// TODO: 함수로 분리
		data, err := file.ReadFile(errorPath)
		if err != nil {
			return RunResult{}, fmt.Errorf("runner: failed to read error file: %w", err)
		}
		runResult.ResultCode = res.ResultCode
		runResult.ErrOutput = string(data)
		fmt.Println(res) // log?
	}

	data, err := file.ReadFile(outputPath)
	if err != nil {
		return RunResult{}, fmt.Errorf("runner: failed to read output file: %w", err)
	}
	runResult.Output = data

	// fmt.Println(res)
	return runResult, nil
}
