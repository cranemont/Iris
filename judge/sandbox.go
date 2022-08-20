package judge

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Sandbox interface {
	Execute(args *SandboxArgs)
}

type SandboxArgs struct {
	ExePath         string
	MaxCpuTime      int
	MaxRealTime     int
	MaxMemory       int
	Args            []string
	InputPath       string
	OutputPath      string
	ErrorPath       string
	LogPath         string
	SeccompRuleName string
	Uid             int
	Gid             int
}

type sandbox struct {
}

func NewSandbox() *sandbox {
	return &sandbox{}
}

func (s *sandbox) isEmptyString(str string) bool {
	if str == "" {
		return true
	}
	return false
}

func (s *sandbox) isEmptyNum(num int) bool {
	if num == 0 {
		return true
	}
	return false
}

func (s *sandbox) Execute(args *SandboxArgs) {
	fmt.Println("Sandbox: func execute")
	exePath := "--exe_path=" + args.ExePath
	maxCpuTime := "--max_cpu_time=" + fmt.Sprint(args.MaxCpuTime)
	maxRealTime := "--max_real_time=" + fmt.Sprint(args.MaxRealTime)
	maxMemory := "--max_memory=" + fmt.Sprint(args.MaxMemory)
	outputPath := "--output_path=out.out"
	errorPath := "--error_path=error.out"

	argsWithFormat := []string{}
	for _, arg := range args.Args {
		argsWithFormat = append(argsWithFormat, "--args="+arg)
	}

	env := "--env=PATH=" + os.Getenv("PATH")
	// envWithFormat := []string{}
	// for _, e := range env {
	// 	envWithFormat = append(envWithFormat, "--env="+e)
	// }

	argSlice := []string{
		exePath, maxCpuTime, maxRealTime, maxMemory, outputPath, errorPath, env, "--uid=0", "--gid=0",
	}
	// argSlice = append(argSlice, append(argsWithFormat, envWithFormat...)...)
	argSlice = append(argSlice, argsWithFormat...)
	// fmt.Print(argSlice)

	cmd := exec.Command("/usr/lib/judger/libjudger.so", argSlice...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
	// stdin, out 연결해서 실행?
}
