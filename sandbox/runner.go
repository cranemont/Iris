package sandbox

import (
	"fmt"
)

type Runner interface {
	Run(dir string, language string, input []byte) (RunResult, error)
}

type runner struct {
	sandbox Sandbox
	config  *LanguageConfig
}

type RunResult struct {
	Signal     int
	ErrorCode  int
	ExitCode   int
	ResultCode int
	Output     []byte
}

func NewRunner(sandbox Sandbox, config *LanguageConfig) *runner {
	return &runner{sandbox, config}
}

func (r *runner) Run(dir string, language string, input []byte) (RunResult, error) {
	fmt.Println("RUN! from runner")

	exePath, err := r.config.MakeExePath(dir, language)
	if err != nil {
		return RunResult{}, err
	}

	//task의 limit으로 주기
	args := ExecArgs{
		ExePath:       exePath,
		MaxCpuTime:    1000,              //task.limit.Time,
		MaxRealTime:   3000,              //task.limit.Time * 3,
		MaxMemory:     256 * 1024 * 1024, //task.limit.Memory,
		MaxStackSize:  128 * 1024 * 1024,
		MaxOutputSize: 10 * 1024 * 1024, // TODO: Testcase 크기 따라서 설정
		OutputPath:    "./run/out.out",
		ErrorPath:     "./run/error.out",
		LogPath:       "./run/log.out",
	}
	result, err := r.sandbox.Execute(args, []byte("input1\ninput2\n"))
	if err != nil {
		return RunResult{}, err
	}
	fmt.Println(result)
	return RunResult{ExitCode: 0}, nil
}

// "command": "{exe_path}",
//         "seccomp_rule": {ProblemIOMode.standard: "c_cpp", ProblemIOMode.file: "c_cpp_file_io"},
//         "env": default_env
