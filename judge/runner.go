package judge

import (
	"fmt"

	"github.com/cranemont/judge-manager/common/dto"
	"github.com/cranemont/judge-manager/judge/config"
)

type Runner interface {
	Run(out chan<- dto.GoResult, task *Task)
}

type runner struct {
	sandbox Sandbox
	config  *config.LanguageConfig
}

type RunResult struct {
	Signal     int
	ErrorCode  int
	ExitCode   int
	ResultCode int
	Output     []byte
}

func NewRunner(sandbox Sandbox, config *config.LanguageConfig) *runner {
	return &runner{sandbox, config}
}

func (r *runner) Run(out chan<- dto.GoResult, task *Task) {
	fmt.Println("RUN! from runner")

	exePath, err := r.config.MakeExePath(task.dir, task.language)
	if err != nil {
		out <- dto.GoResult{Err: err, Data: RunResult{}}
		return
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
	r.sandbox.Execute(args, []byte("input1\ninput2\n"))
	out <- dto.GoResult{Data: RunResult{ExitCode: 0}}
	// 채널로 결과반환
}

// "command": "{exe_path}",
//         "seccomp_rule": {ProblemIOMode.standard: "c_cpp", ProblemIOMode.file: "c_cpp_file_io"},
//         "env": default_env
