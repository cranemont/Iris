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

func (s *sandbox) Execute(args *SandboxArgs) {
	fmt.Println("Sandbox: func execute")

	exePath := "--exe_path=" + args.ExePath
	maxCpuTime := "--max_cpu_time=" + fmt.Sprint(args.MaxCpuTime)
	maxRealTime := "--max_real_time=" + fmt.Sprint(args.MaxRealTime)
	maxMemory := "--max_memory=" + fmt.Sprint(args.MaxMemory)
	// outputPath := "--output_path=./compiler.out"
	// errorPath := "--error_path=./compiler.out"

	argsWithFormat := []string{}
	for _, arg := range args.Args {
		argsWithFormat = append(argsWithFormat, "--args="+arg)
	}

	argSlice := []string{
		exePath, maxCpuTime, maxRealTime, maxMemory,
	}
	argSlice = append(argSlice, (argsWithFormat)...)
	fmt.Print(argSlice)

	cmd := exec.Command("/usr/lib/judger/libjudger.so", argSlice...)
	// cmd := exec.Command("/usr/bin/gcc", args.Args[5])
	cmd.Env = os.Environ()
	// cmd.Stdin = strings.NewReader("input from judger")
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
	// out, _ := cmd.Output()
	// err := cmd.Run()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(out))
	// stdin, out 연결해서 실행?
}
