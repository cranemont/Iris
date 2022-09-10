package sandbox

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cranemont/judge-manager/common/file"
)

const (
	C    = "c"
	CPP  = "cpp"
	JAVA = "java"
	// PYTHON2 = "py2"
	PYTHON = "py3"
)

func GetConfig(language string) (LangConfig, error) {
	switch language {
	case C:
		return cConfig, nil
	case CPP:
		return cppConfig, nil
	case JAVA:
		return javaConfig, nil
	// case PYTHON2:
	// 	return py2Config, nil
	case PYTHON:
		return pyConfig, nil
	}
	return LangConfig{}, fmt.Errorf("unsupported language: %s", language)
}

type LangConfig struct {
	Language              string
	SrcName               string
	ExeName               string
	MaxCompileCpuTime     int
	MaxCompileRealTime    int
	MaxCompileMemory      int
	CompilerPath          string
	CompileArgs           string
	RunCommand            string
	RunArgs               string
	SeccompRule           string
	SeccompRuleFileIO     string
	MemoeryLimitCheckOnly bool
	env                   []string
}

func (c LangConfig) SrcPath(dir string) string {
	return file.MakeFilePath(dir, c.SrcName).String()
}

func (c LangConfig) CompileOutputPath(dir string) string {
	return file.MakeFilePath(dir, CompileOutFile).String()
}

func (c LangConfig) RunOutputPath(dir string, order int) string {
	return file.MakeFilePath(dir, strconv.Itoa(order)+".out").String()
}

func (c LangConfig) RunErrPath(dir string, order int) string {
	return file.MakeFilePath(dir, strconv.Itoa(order)+".error").String()
}

func (c LangConfig) ToCompileExecArgs(dir string) ExecArgs {

	outputPath := file.MakeFilePath(dir, CompileOutFile).String()
	srcPath := file.MakeFilePath(dir, c.SrcName).String()
	exePath := file.MakeFilePath(dir, c.ExeName).String()
	exeDir := file.MakeFilePath(dir, "").String()

	args := strings.Replace(c.CompileArgs, "{srcPath}", srcPath, 1)
	args = strings.Replace(args, "{exePath}", exePath, 1)
	args = strings.Replace(args, "{exeDir}", exeDir, 1)

	return ExecArgs{
		ExePath:      c.CompilerPath,
		MaxCpuTime:   c.MaxCompileCpuTime,
		MaxRealTime:  c.MaxCompileRealTime,
		MaxMemory:    c.MaxCompileMemory,
		MaxStackSize: 128 * 1024 * 1024,
		// FIXME: testcase크기 따라서 설정하거나, 그냥 바로 stdout 읽어오거나
		MaxOutputSize: 20 * 1024 * 1024,
		OutputPath:    outputPath,
		ErrorPath:     outputPath,
		LogPath:       CompileLogPath,
		Args:          validateArgs(args),
	}
}

func validateArgs(args string) []string {
	if args != "" {
		return strings.Split(args, " ")
	}
	return nil
}

func (c LangConfig) ToRunExecArgs(dir string, order int, limit Limit, fileIo bool) ExecArgs {
	exePath := file.MakeFilePath(dir, c.ExeName).String()
	exeDir := file.MakeFilePath(dir, "").String()
	outputPath := file.MakeFilePath(dir, strconv.Itoa(order)+".out").String()
	errorPath := file.MakeFilePath(dir, strconv.Itoa(order)+".error").String()

	// run args 설정
	args := strings.Replace(c.RunArgs, "{maxMemory}", strconv.Itoa(limit.Memory), 1)
	args = strings.Replace(args, "{exePath}", exePath, 1)
	args = strings.Replace(args, "{exeDir}", exeDir, 1)

	maxMemory := limit.Memory
	if c.Language == JAVA {
		maxMemory = -1
	}

	return ExecArgs{
		ExePath:      strings.Replace(c.RunCommand, "{exePath}", exePath, 1),
		MaxCpuTime:   limit.CpuTime,
		MaxRealTime:  limit.RealTime,
		MaxMemory:    maxMemory,
		MaxStackSize: 128 * 1024 * 1024,
		// FIXME: testcase크기 따라서 설정하거나, 그냥 바로 stdout 읽어오거나 -> 이러면 order 필요없음
		MaxOutputSize: 10 * 1024 * 1024,
		// file에 쓰는거랑 stdout이랑 크게 차이 안남
		// https://stackoverflow.com/questions/29700478/redirecting-of-stdout-in-bash-vs-writing-to-file-in-c-with-fprintf-speed
		OutputPath:      outputPath,
		ErrorPath:       errorPath, // byte buffer로
		LogPath:         RunLogPath,
		SeccompRuleName: c.SeccompRule,
		Args:            validateArgs(args),
	}
}

type Limit struct {
	CpuTime  int
	RealTime int
	Memory   int
}

// srcPath, exePath는 base dir + task dir
var defaultEnv = []string{"LANG=en_US.UTF-8", "LANGUAGE=en_US:en", "LC_ALL=en_US.UTF-8"}
var cConfig = LangConfig{
	Language:           C,
	SrcName:            "main.c",
	ExeName:            "main",
	MaxCompileCpuTime:  3000,              // compile
	MaxCompileRealTime: 10000,             // compile
	MaxCompileMemory:   256 * 1024 * 1024, // compile
	CompilerPath:       "/usr/bin/gcc",
	CompileArgs: "-DONLINE_JUDGE " +
		"-O2 -w -fmax-errors=3 -std=c11 " +
		"{srcPath} -lm -o {exePath}",
	RunCommand:            "{exePath}",
	RunArgs:               "",
	SeccompRule:           "c_cpp",
	SeccompRuleFileIO:     "c_cpp_file_io",
	MemoeryLimitCheckOnly: false,
	env:                   defaultEnv,
}

var cppConfig = LangConfig{
	Language:           CPP,
	SrcName:            "main.cpp",
	ExeName:            "main",
	MaxCompileCpuTime:  10000,
	MaxCompileRealTime: 20000,
	MaxCompileMemory:   1024 * 1024 * 1024,
	CompilerPath:       "/usr/bin/g++",
	CompileArgs: "-DONLINE_JUDGE " +
		"-O2 -w -fmax-errors=3 " +
		"-std=c++14 {srcPath} -lm -o {exePath}",
	RunCommand:            "{exePath}",
	RunArgs:               "",
	SeccompRule:           "c_cpp",
	SeccompRuleFileIO:     "c_cpp_file_io",
	MemoeryLimitCheckOnly: false,
	env:                   defaultEnv,
}

var javaConfig = LangConfig{
	Language:           JAVA,
	SrcName:            "Main.java",
	ExeName:            "Main",
	MaxCompileCpuTime:  5000,
	MaxCompileRealTime: 10000,
	MaxCompileMemory:   -1,
	CompilerPath:       "/usr/bin/javac",
	CompileArgs:        "{srcPath} -d {exeDir} -encoding UTF8",
	RunCommand:         "/usr/bin/java",
	RunArgs: "-cp {exeDir} " +
		"-XX:MaxRAM={maxMemory}k " +
		"-Djava.security.manager " +
		"-Dfile.encoding=UTF-8 " +
		"-Djava.security.policy==/etc/java_policy " +
		"-Djava.awt.headless=true " +
		"Main",
	SeccompRule:           "",
	MemoeryLimitCheckOnly: true,
	env:                   defaultEnv,
}

// var py2Config = LangConfig{
// 	Language:              PYTHON2,
// 	SrcName:               "solution.py",
// 	ExeName:               "solution.pyc",
// 	MaxCompileCpuTime:     3000,
// 	MaxCompileRealTime:    10000,
// 	MaxCompileMemory:      1024 * 1024 * 1024,
// 	CompilerPath:          "/usr/bin/python",
// 	CompileArgs:           "-m py_compile {srcPath}",
// 	RunCommand:            "/usr/bin/python",
// 	RunArgs:               "{exePath}",
// 	SeccompRule:           "general",
// 	MemoeryLimitCheckOnly: false,
// 	env:                   defaultEnv,
// }

var pyConfig = LangConfig{
	Language:              PYTHON,
	SrcName:               "solution.py",
	ExeName:               "__pycache__/solution.cpython-38.pyc", // TODO: 파이썬 버전 확인
	MaxCompileCpuTime:     3000,
	MaxCompileRealTime:    10000,
	MaxCompileMemory:      128 * 1024 * 1024,
	CompilerPath:          "/usr/bin/python3",
	CompileArgs:           "-m py_compile {srcPath}",
	RunCommand:            "/usr/bin/python3",
	RunArgs:               "{exePath}",
	SeccompRule:           "general",
	MemoeryLimitCheckOnly: false,
	env:                   append(defaultEnv, "PYTHONIOENCODING=utf-8"),
}
