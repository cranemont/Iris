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
}

func NewRunner(sandbox Sandbox, config *config.LanguageConfig) *runner {
	return &runner{sandbox, config}
}

func (r *runner) Run(out chan<- dto.GoResult, task *Task) {
	fmt.Println("RUN! from runner")

	options, err := r.config.Get(task.language) // 이게 된다고? private 아닌가? GetLanguage 가 필요없어?
	if err != nil {
		err := fmt.Errorf("failed to get language config: %s", err)
		out <- dto.GoResult{Err: err, Data: RunResult{}}
		return
	}

	exePath, err := r.config.GetExePath(task.dir, task.language)
	if err != nil {
		err := fmt.Errorf("failed to get language config: %s", err)
		out <- dto.GoResult{Err: err, Data: RunResult{}}
		return
	}

	args := SandboxArgs{
		ExePath:     exePath,
		MaxCpuTime:  options.MaxCpuTime,
		MaxRealTime: options.MaxRealTime,
		MaxMemory:   options.MaxMemory,
	}
	r.sandbox.Run(&args)
	out <- dto.GoResult{Data: RunResult{ExitCode: 0}}
	// 채널로 결과반환
}

// "command": "{exe_path}",
//         "seccomp_rule": {ProblemIOMode.standard: "c_cpp", ProblemIOMode.file: "c_cpp_file_io"},
//         "env": default_env
