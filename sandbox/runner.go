package sandbox

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cranemont/judge-manager/file"
)

type Runner interface {
	Run(dto RunRequest, input []byte) (RunResult, error)
}

type runner struct {
	config *LanguageConfig
}

type RunResult struct {
	Order      int
	ResultCode int
	ErrOutput  string // []byte?
	ExecResult string
	Output     []byte
}

type RunRequest struct {
	Order       int
	Dir         string
	Language    string
	TimeLimit   int
	MemoryLimit int
}

func NewRunner(config *LanguageConfig) *runner {
	return &runner{config}
}

func (r *runner) Run(dto RunRequest, input []byte) (RunResult, error) {
	fmt.Println("RUN! from runner")
	dir := dto.Dir
	id := dto.Order
	language := dto.Language
	timeLimit := dto.TimeLimit
	memoryLimit := dto.MemoryLimit

	exePath, err := r.config.MakeExePath(dir, language)
	if err != nil {
		return RunResult{}, err
	}

	outputPath := file.MakeFilePath(dir, strconv.Itoa(id)+".out").String()
	errorPath := file.MakeFilePath(dir, strconv.Itoa(id)+".error").String()
	//task의 limit으로 주기
	res, err := Exec(
		ExecArgs{
			ExePath:       exePath,
			MaxCpuTime:    timeLimit,     //task.limit.Time,
			MaxRealTime:   timeLimit * 3, //task.limit.Time * 3, 다른 task들에 영향받을 수 있기 때문
			MaxMemory:     memoryLimit,   //task.limit.Memory,
			MaxStackSize:  128 * 1024 * 1024,
			MaxOutputSize: 10 * 1024 * 1024, // TODO: Testcase 크기 따라서 설정
			// file에 쓰는거랑 stdout이랑 크게 차이 안남
			// https://stackoverflow.com/questions/29700478/redirecting-of-stdout-in-bash-vs-writing-to-file-in-c-with-fprintf-speed
			OutputPath: outputPath,
			ErrorPath:  errorPath, // byte buffer로
			LogPath:    "./log/run/log.out",
		}, input,
	)
	if err != nil {
		return RunResult{}, err
	}

	runResult := RunResult{Order: id, ResultCode: SUCCESS}
	if res.ResultCode != SUCCESS {
		// TODO: 함수로 분리
		sandboxResult, err := json.Marshal(res)
		if err != nil {
			return RunResult{}, fmt.Errorf("invalid result format: %w", err)
		}
		data, err := file.ReadFile(errorPath)
		if err != nil {
			return RunResult{}, fmt.Errorf("failed to read output file: %w", err)
		}
		runResult.ResultCode = res.ResultCode
		runResult.ExecResult = string(sandboxResult) // 필요한가?
		runResult.ErrOutput = string(data)
		fmt.Println(runResult)
	}

	data, err := file.ReadFile(outputPath)
	if err != nil {
		return RunResult{}, err
	}
	runResult.Output = data

	fmt.Println(res)
	return runResult, nil
}

// "command": "{exe_path}",
//         "seccomp_rule": {ProblemIOMode.standard: "c_cpp", ProblemIOMode.file: "c_cpp_file_io"},
//         "env": default_env
