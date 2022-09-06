package sandbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func Exec(args ExecArgs, input []byte) (SandboxResult, error) {
	// fmt.Println("input: ", args)
	argSlice := makeExecArgs(args)
	env := "--env=PATH=" + os.Getenv("PATH")
	argSlice = append(argSlice, env)

	// fmt.Println(argSlice)
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
		log.Println(fmt.Sprint(err) + ": " + stderr.String())
		return SandboxResult{}, fmt.Errorf("execution failed: %w: %s", err, stderr.String())
	}

	res := SandboxResult{}

	json.Unmarshal(stdout.Bytes(), &res)
	fmt.Println("Result: ", stdout.String()) // on debug

	return res, nil
}

type SandboxResult struct {
	CpuTime    int `json:"cpuTime"`
	RealTime   int `json:"realTime"`
	Memory     int `json:"memory"`
	Signal     int `json:"signal"`
	ErrorCode  int `json:"exitCode"`
	ExitCode   int `json:"errorCode"`
	ResultCode int `json:"resultCode"`
}

type ExecArgs struct {
	MaxCpuTime           int
	MaxRealTime          int
	MaxMemory            int
	MaxStackSize         int
	MaxOutputSize        int
	ExePath              string
	InputPath            string
	OutputPath           string
	ErrorPath            string
	LogPath              string
	SeccompRuleName      string
	MemoryLimitCheckOnly bool
	Args                 []string
	Env                  []string
	Uid                  int
	Gid                  int
}

type formatString struct {
	MaxCpuTime           string
	MaxRealTime          string
	MaxMemory            string
	MaxStackSize         string
	MaxOutputSize        string
	ExePath              string
	InputPath            string
	OutputPath           string
	ErrorPath            string
	LogPath              string
	Args                 string
	Env                  string
	SeccompRuleName      string
	MemoryLimitCheckOnly string
	Uid                  string
	Gid                  string
}

var format = formatString{
	MaxCpuTime:           "--max_cpu_time=",
	MaxRealTime:          "--max_real_time=",
	MaxMemory:            "--max_memory=",
	MaxStackSize:         "--max_stack=",
	MaxOutputSize:        "--max_output_size=",
	ExePath:              "--exe_path=",
	InputPath:            "--input_path=",
	OutputPath:           "--output_path=",
	ErrorPath:            "--error_path=",
	LogPath:              "--log_path=",
	Args:                 "--args=",
	Env:                  "--env=",
	SeccompRuleName:      "--seccomp_rule_name=",
	MemoryLimitCheckOnly: "--memory_limit_check_only=",
	Uid:                  "--uid=",
	Gid:                  "--gid=",
}

// methods below is for the libjudger specific
func makeExecArgs(data ExecArgs) []string {
	argSlice := []string{}
	if !isEmptyInt(data.MaxCpuTime) {
		argSlice = concatIntArgs(argSlice, format.MaxCpuTime, data.MaxCpuTime)
	}
	if !isEmptyInt(data.MaxRealTime) {
		argSlice = concatIntArgs(argSlice, format.MaxRealTime, data.MaxRealTime)
	}
	if !isEmptyInt(data.MaxMemory) {
		argSlice = concatIntArgs(argSlice, format.MaxMemory, data.MaxMemory)
	}
	if !isEmptyInt(data.MaxStackSize) {
		argSlice = concatIntArgs(argSlice, format.MaxStackSize, data.MaxStackSize)
	}
	if !isEmptyInt(data.MaxOutputSize) {
		argSlice = concatIntArgs(argSlice, format.MaxOutputSize, data.MaxOutputSize)
	}
	if data.Uid >= 0 && data.Uid < 65534 {
		// FIXME: set default uid
		argSlice = concatIntArgs(argSlice, format.Uid, data.Uid)
	}
	if data.Gid >= 0 && data.Gid < 65534 {
		// FIXME: set default Gid
		argSlice = concatIntArgs(argSlice, format.Gid, data.Gid)
	}
	if !isEmptyString(data.ExePath) {
		argSlice = concatStringArgs(argSlice, format.ExePath, data.ExePath)
	}
	if !isEmptyString(data.InputPath) {
		argSlice = concatStringArgs(argSlice, format.InputPath, data.InputPath)
	}
	if !isEmptyString(data.OutputPath) {
		argSlice = concatStringArgs(argSlice, format.OutputPath, data.OutputPath)
	}
	if !isEmptyString(data.ErrorPath) {
		argSlice = concatStringArgs(argSlice, format.ErrorPath, data.ErrorPath)
	}
	if !isEmptyString(data.LogPath) {
		argSlice = concatStringArgs(argSlice, format.LogPath, data.LogPath)
	}
	if !isEmptyString(data.SeccompRuleName) {
		argSlice = concatStringArgs(argSlice, format.SeccompRuleName, data.SeccompRuleName)
	}
	// TODO: MemoryLimitCheckOnly 추가
	if !isEmptySlice(data.Args) {
		argSlice = concatSliceArgs(argSlice, format.Args, data.Args)
	}
	if !isEmptySlice(data.Env) {
		argSlice = concatSliceArgs(argSlice, format.Env, data.Env)
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
