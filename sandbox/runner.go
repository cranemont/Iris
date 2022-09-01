package sandbox

import (
	"fmt"
	"strconv"

	"github.com/cranemont/judge-manager/file"
)

type Runner interface {
	Run(dir string, id int, language string, input []byte) (RunResult, error)
}

type runner struct {
	sandbox Sandbox
	config  *LanguageConfig
}

type RunResult struct {
	Id         int
	Signal     int
	ErrorCode  int
	ExitCode   int
	ResultCode int
	Output     []byte
}

func NewRunner(sandbox Sandbox, config *LanguageConfig) *runner {
	return &runner{sandbox, config}
}

func (r *runner) Run(dir string, id int, language string, input []byte) (RunResult, error) {
	fmt.Println("RUN! from runner")

	exePath, err := r.config.MakeExePath(dir, language)
	if err != nil {
		return RunResult{}, err
	}

	outputPath := MakeFilePath(dir, strconv.Itoa(id)+".out").String()
	errorPath := MakeFilePath(dir, strconv.Itoa(id)+".error").String()
	//task의 limit으로 주기
	args := ExecArgs{
		ExePath:       exePath,
		MaxCpuTime:    1000,              //task.limit.Time,
		MaxRealTime:   3000,              //task.limit.Time * 3,
		MaxMemory:     256 * 1024 * 1024, //task.limit.Memory,
		MaxStackSize:  128 * 1024 * 1024,
		MaxOutputSize: 10 * 1024 * 1024, // TODO: Testcase 크기 따라서 설정
		// file에 쓰는거랑 stdout이랑 크게 차이 안남
		// https://stackoverflow.com/questions/29700478/redirecting-of-stdout-in-bash-vs-writing-to-file-in-c-with-fprintf-speed
		OutputPath: outputPath,
		ErrorPath:  errorPath, // byte buffer로
		LogPath:    "./log/run/log.out",
	}
	result, err := r.sandbox.Execute(args, input)
	if err != nil {
		return RunResult{}, err
	}

	data, err := file.ReadFile(outputPath)
	if err != nil {
		return RunResult{}, err
	}

	fmt.Println(result)
	return RunResult{Id: id, ExitCode: 0, Output: data}, nil
}

// "command": "{exe_path}",
//         "seccomp_rule": {ProblemIOMode.standard: "c_cpp", ProblemIOMode.file: "c_cpp_file_io"},
//         "env": default_env
