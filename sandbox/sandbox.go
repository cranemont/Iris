package sandbox

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type Result struct {
	Code       int
	Signal     int
	ErrorCode  int
	ExitCode   int
	ResultCode int
	Err        string
	Output     []byte
}

type Sandbox interface {
	Execute(args ExecArgs, input []byte) (Result, error)
}

type ExecArgs struct {
	MaxCpuTime      int
	MaxRealTime     int
	MaxMemory       int
	MaxStackSize    int
	MaxOutputSize   int
	ExePath         string
	InputPath       string
	OutputPath      string
	ErrorPath       string
	LogPath         string
	SeccompRuleName string
	Args            []string
	Env             []string
	Uid             int
	Gid             int
}

type formatString struct {
	MaxCpuTime      string
	MaxRealTime     string
	MaxMemory       string
	MaxStackSize    string
	MaxOutputSize   string
	ExePath         string
	InputPath       string
	OutputPath      string
	ErrorPath       string
	LogPath         string
	Args            string
	Env             string
	SeccompRuleName string
	Uid             string
	Gid             string
}

type sandbox struct {
	format formatString
}

func NewSandbox() *sandbox {
	format := formatString{
		MaxCpuTime:      "--max_cpu_time=",
		MaxRealTime:     "--max_real_time=",
		MaxMemory:       "--max_memory=",
		MaxStackSize:    "--max_stack=",
		MaxOutputSize:   "--max_output_size=",
		ExePath:         "--exe_path=",
		InputPath:       "--input_path=",
		OutputPath:      "--output_path=",
		ErrorPath:       "--error_path=",
		LogPath:         "--log_path=",
		Args:            "--args=",
		Env:             "--env=",
		SeccompRuleName: "--seccomp_rule_name=",
		Uid:             "--uid=",
		Gid:             "--gid=",
	}
	return &sandbox{format}
}

func (s *sandbox) Execute(args ExecArgs, input []byte) (Result, error) {
	argSlice := s.makeExecArgs(args)
	env := "--env=PATH=" + os.Getenv("PATH")
	argSlice = append(argSlice, env)

	cmd := exec.Command("/usr/lib/judger/libjudger.so", argSlice...)

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdin = &stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	stdin.Write(input)

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return Result{}, fmt.Errorf("execution fail: %w: %s", err, stderr.String())
	}
	fmt.Println("Result: " + stdout.String())
	// fmt.Println(stdout.Len())
	return Result{Code: 0}, nil
}

// methods below is for the libjudger specific
func (s *sandbox) makeExecArgs(data ExecArgs) []string {
	argSlice := []string{}
	if !isEmptyInt(data.MaxCpuTime) {
		argSlice = concatIntArgs(argSlice, s.format.MaxCpuTime, data.MaxCpuTime)
	}
	if !isEmptyInt(data.MaxRealTime) {
		argSlice = concatIntArgs(argSlice, s.format.MaxRealTime, data.MaxRealTime)
	}
	if !isEmptyInt(data.MaxMemory) {
		argSlice = concatIntArgs(argSlice, s.format.MaxMemory, data.MaxMemory)
	}
	if !isEmptyInt(data.MaxStackSize) {
		argSlice = concatIntArgs(argSlice, s.format.MaxStackSize, data.MaxStackSize)
	}
	if !isEmptyInt(data.MaxOutputSize) {
		argSlice = concatIntArgs(argSlice, s.format.MaxOutputSize, data.MaxOutputSize)
	}
	if data.Uid >= 0 && data.Uid < 65534 {
		argSlice = concatIntArgs(argSlice, s.format.Uid, data.Uid)
	}
	if data.Uid >= 0 && data.Uid < 65534 {
		argSlice = concatIntArgs(argSlice, s.format.Gid, data.Gid)
	}
	if !isEmptyString(data.ExePath) {
		argSlice = concatStringArgs(argSlice, s.format.ExePath, data.ExePath)
	}
	if !isEmptyString(data.InputPath) {
		argSlice = concatStringArgs(argSlice, s.format.InputPath, data.InputPath)
	}
	if !isEmptyString(data.OutputPath) {
		argSlice = concatStringArgs(argSlice, s.format.OutputPath, data.OutputPath)
	}
	if !isEmptyString(data.ErrorPath) {
		argSlice = concatStringArgs(argSlice, s.format.ErrorPath, data.ErrorPath)
	}
	if !isEmptyString(data.LogPath) {
		argSlice = concatStringArgs(argSlice, s.format.LogPath, data.LogPath)
	}
	if !isEmptyString(data.SeccompRuleName) {
		argSlice = concatStringArgs(argSlice, s.format.SeccompRuleName, data.SeccompRuleName)
	}
	if !isEmptySlice(data.Args) {
		argSlice = concatSliceArgs(argSlice, s.format.Args, data.Args)
	}
	if !isEmptySlice(data.Env) {
		argSlice = concatSliceArgs(argSlice, s.format.Env, data.Env)
	}
	return argSlice
}

func isEmptyString(str string) bool {
	return str == ""
}

func isEmptyInt(num int) bool {
	return num == 0
}

func isEmptySlice(slice []string) bool {
	return slice == nil
}

func concatStringArgs(argSlice []string, format string, arg string) []string {
	var b bytes.Buffer
	b.WriteString(format)
	b.WriteString(arg)
	return append(argSlice, b.String())
}

func concatIntArgs(argSlice []string, format string, arg int) []string {
	var b bytes.Buffer
	b.WriteString(format)
	b.WriteString(strconv.Itoa(arg))
	return append(argSlice, b.String())
}

func concatSliceArgs(argSlice []string, format string, arg []string) []string {
	var b bytes.Buffer
	for _, arg := range arg {
		b.WriteString(format)
		b.WriteString(arg)
		argSlice = append(argSlice, b.String())
		b.Reset()
	}
	return argSlice
}
