package sandbox

import (
	"fmt"
)

type Runner interface {
	Run(dir string, id string, language string, input []byte) (RunResult, error)
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

func (r *runner) Run(dir string, id string, language string, input []byte) (RunResult, error) {
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
		// file에 쓰는거랑 stdout이랑 크게 차이 안남
		// https://stackoverflow.com/questions/29700478/redirecting-of-stdout-in-bash-vs-writing-to-file-in-c-with-fprintf-speed
		OutputPath: MakeFilePath(dir, id+".out").String(),
		ErrorPath:  "./run/error.out", //compile은 되는데 run은 안되는 상황에서 error가 덮어씌워지는지?
		LogPath:    "./run/log.out",
	}
	result, err := r.sandbox.Execute(args, input)
	if err != nil {
		return RunResult{}, err
	}

	fmt.Println(result)
	return RunResult{ExitCode: 0}, nil
}

// "command": "{exe_path}",
//         "seccomp_rule": {ProblemIOMode.standard: "c_cpp", ProblemIOMode.file: "c_cpp_file_io"},
//         "env": default_env
