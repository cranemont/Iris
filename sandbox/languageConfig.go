package sandbox

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cranemont/judge-manager/constants"
)

// 언어별 설정 관리..
// 각 DTO로 변환 과정 수행

type config struct {
	SrcName      string
	ExeName      string
	MaxCpuTime   int
	MaxRealTime  int
	MaxMemory    int
	CompilerPath string
	Args         string
}
type LanguageConfig struct {
	// Get(language) -> constants패키지에서 설정값 가져옴
	configMap map[string]config
}

func (l *LanguageConfig) Init() {
	l.configMap = make(map[string]config)
	// srcPath, exePath는 base dir + task dir
	l.configMap["c"] = config{
		SrcName:      "main.c",
		ExeName:      "main",
		MaxCpuTime:   3000,
		MaxRealTime:  10000,
		MaxMemory:    256 * 1024 * 1024,
		CompilerPath: "/usr/bin/gcc",
		Args:         "-DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c11 {srcPath} -lm -o {exePath}",
	}
	l.configMap["cpp"] = config{
		SrcName:      "main.cpp",
		ExeName:      "main",
		MaxCpuTime:   10000,
		MaxRealTime:  20000,
		MaxMemory:    1024 * 1024 * 1024,
		CompilerPath: "/usr/bin/g++",
		Args:         "-DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c++14 {src_path} -lm -o {exe_path}",
	}
}

// enum으로 바꾸기
func (l *LanguageConfig) Get(language string) (config, error) {
	if val, ok := l.configMap[language]; ok {
		return val, nil
	}
	return config{}, fmt.Errorf("unsupported language: %s", language)
}

func (l *LanguageConfig) MakeSrcPath(dir string, language string) (string, error) {
	conf, err := l.Get(language)
	if err != nil {
		return "", fmt.Errorf("failed to make srcpath: %w", err)
	}
	return MakeFilePath(dir, conf.SrcName).String(), nil
}

func (l *LanguageConfig) MakeExePath(dir string, language string) (string, error) {
	conf, err := l.Get(language)
	if err != nil {
		return "", fmt.Errorf("failed to make exepath: %w", err)
	}
	return MakeFilePath(dir, conf.ExeName).String(), nil
}

func (l *LanguageConfig) MakeArgSlice(srcPath string, exePath string, language string) ([]string, error) {

	conf, err := l.Get(language)
	if err != nil {
		return nil, fmt.Errorf("failed to make argslice: %w", err)
	}
	args := strings.Replace(conf.Args, "{srcPath}", srcPath, 1)
	args = strings.Replace(args, "{exePath}", exePath, 1)
	return strings.Split(args, " "), nil
}

func MakeFilePath(dir string, fileName string) *bytes.Buffer {
	var b bytes.Buffer
	b.WriteString(constants.BASE_DIR)
	b.WriteString("/")
	b.WriteString(dir)
	b.WriteString("/")
	return &b
}

// "compile": {
//         "src_name": "main.c",
//         "exe_name": "main",
//         "max_cpu_time": 3000,
//         "max_real_time": 10000,
//         "max_memory": 256 * 1024 * 1024,
//         "command": "/usr/bin/gcc -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c11 {src_path} -lm -o {exe_path}",
//     },
//     "run": {
//         "command": "{exe_path}",
//         "seccomp_rule": {ProblemIOMode.standard: "c_cpp", ProblemIOMode.file: "c_cpp_file_io"},
//         "env": default_env
//     }
