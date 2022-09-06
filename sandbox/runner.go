package sandbox

import (
	"fmt"
	"strconv"

	"github.com/cranemont/judge-manager/file"
)

type Runner interface {
	Run(dto RunRequest, input []byte) (RunResult, error)
}

type runner struct {
}

type RunResult struct {
	Order      int
	ResultCode int
	ErrOutput  string // []byte?
	CpuTime    int
	RealTime   int
	Memory     int
	Signal     int
	ErrorCode  int
	ExitCode   int
	Output     []byte
}

type RunRequest struct {
	Order       int
	Dir         string
	Language    string
	TimeLimit   int
	MemoryLimit int
}

func NewRunner() *runner {
	return &runner{}
}

func (r *runner) Run(dto RunRequest, input []byte) (RunResult, error) {
	fmt.Println("RUN! from runner")
	dir := dto.Dir
	order := dto.Order
	language := dto.Language
	timeLimit := dto.TimeLimit
	memoryLimit := dto.MemoryLimit

	exePath, err := r.config.MakeExePath(dir, language)
	if err != nil {
		return RunResult{}, err
	}
	// argSlice, err := r.config.MakeRunArgSlice(srcPath, exePath, language)
	// if err != nil {
	// 	return CompileResult{}, err
	// }

	outputPath := file.MakeFilePath(dir, strconv.Itoa(order)+".out").String()
	errorPath := file.MakeFilePath(dir, strconv.Itoa(order)+".error").String()
	//task의 limit으로 주기
	res, err := Exec(
		ExecArgs{
			ExePath:       exePath,
			MaxCpuTime:    timeLimit,     //task.limit.Time,
			MaxRealTime:   timeLimit * 3, //task.limit.Time * 3, 다른 task들에 영향받을 수 있기 때문
			MaxMemory:     memoryLimit,   //task.limit.Memory,
			MaxStackSize:  128 * 1024 * 1024,
			MaxOutputSize: 10 * 1024 * 1024, // TODO: Testcase 크기 따라서 설정
			// file에 쓰는거랑 stdout이랑 크게 차이 안남
			// https://stackoverflow.com/questions/29700478/redirecting-of-stdout-in-bash-vs-writing-to-file-in-c-with-fprintf-speed
			OutputPath: outputPath,
			ErrorPath:  errorPath, // byte buffer로
			LogPath:    RunLogPath,
		}, input,
	)
	if err != nil {
		return RunResult{}, fmt.Errorf("runner: Run failed: %w", err)
	}

	runResult := RunResult{
		Order:      order,
		ResultCode: SUCCESS,
		CpuTime:    res.CpuTime,
		RealTime:   res.RealTime,
		Memory:     res.Memory,
		Signal:     res.Signal,
		ErrorCode:  res.ErrorCode,
		ExitCode:   res.ExitCode,
	}

	if res.ResultCode != SUCCESS {
		// TODO: 함수로 분리
		data, err := file.ReadFile(errorPath)
		if err != nil {
			return RunResult{}, fmt.Errorf("runner: failed to read error file: %w", err)
		}
		runResult.ResultCode = res.ResultCode
		runResult.ErrOutput = string(data)
		fmt.Println(res) // log?
	}

	data, err := file.ReadFile(outputPath)
	if err != nil {
		return RunResult{}, fmt.Errorf("runner: failed to read output file: %w", err)
	}
	runResult.Output = data

	// fmt.Println(res)
	return runResult, nil
}
