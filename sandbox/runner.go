package sandbox

import (
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
	Id         int
	Signal     int
	ErrorCode  int
	ExitCode   int
	ResultCode int
	Output     []byte
}

type RunRequest struct {
	Id          int
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
	id := dto.Id
	language := dto.Language
	timeLimit := dto.TimeLimit
	memoryLimit := dto.MemoryLimit

	exePath, err := r.config.MakeExePath(dir, language)
	if err != nil {
		return RunResult{}, err
	}

	outputPath := file.MakeFilePath(dto.Dir, strconv.Itoa(id)+".out").String()
	errorPath := file.MakeFilePath(dto.Dir, strconv.Itoa(id)+".error").String()
	//task의 limit으로 주기
	result, err := Exec(
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

	data, err := file.ReadFile(outputPath)
	if err != nil {
		return RunResult{}, err
	}

	fmt.Println(result)
	return RunResult{Id: dto.Id, ExitCode: 0, Output: data}, nil
}

// "command": "{exe_path}",
//         "seccomp_rule": {ProblemIOMode.standard: "c_cpp", ProblemIOMode.file: "c_cpp_file_io"},
//         "env": default_env
