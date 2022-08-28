package judge

import (
	"fmt"

	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/judge/config"
)

type Runner interface {
	Run(task *Task, out chan<- string)
}

type runner struct {
	sandbox Sandbox
	option  *config.CompileOption
}

func NewRunner(sandbox Sandbox, option *config.CompileOption) *runner {
	return &runner{sandbox, option}
}

func (r *runner) Run(task *Task, out chan<- string) {
	fmt.Println("RUN! from runner")

	options := r.option.Get(task.language) // 이게 된다고? private 아닌가? GetLanguage 가 필요없어?
	exePath := constants.BASE_DIR + "/" + task.GetDir() + "/" + options.ExeName

	args := SandboxArgs{
		ExePath:     exePath,
		MaxCpuTime:  options.MaxCpuTime,
		MaxRealTime: options.MaxRealTime,
		MaxMemory:   options.MaxMemory,
	}
	r.sandbox.Run(&args)
	dir := task.GetDir()
	out <- "task " + dir + " running done"
	// 채널로 결과반환
}

// "command": "{exe_path}",
//         "seccomp_rule": {ProblemIOMode.standard: "c_cpp", ProblemIOMode.file: "c_cpp_file_io"},
//         "env": default_env
