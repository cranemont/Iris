package sandbox

import (
	"fmt"
	"strings"

	"github.com/cranemont/judge-manager/file"
)

// 언어별 설정 관리..
// 각 DTO로 변환 과정 수행

type config struct {
	SrcName           string
	ExeName           string
	MaxCpuTime        int
	MaxRealTime       int
	MaxMemory         int
	CompilerPath      string
	CompileArgs       string
	RunCommand        string
	SeccompRule       string
	SeccompRuleFileIO string
}
type LanguageConfig struct {
	configMap map[string]config
}

func (l *LanguageConfig) Init() {
	l.configMap = make(map[string]config)
	// srcPath, exePath는 base dir + task dir
	l.configMap["c"] = config{
		SrcName:           "main.c",
		ExeName:           "main",
		MaxCpuTime:        3000,              // compile
		MaxRealTime:       10000,             // compile
		MaxMemory:         256 * 1024 * 1024, // compile
		CompilerPath:      "/usr/bin/gcc",
		CompileArgs:       "-DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c11 {srcPath} -lm -o {exePath}",
		RunCommand:        "{exePath}",
		SeccompRule:       "c_cpp",
		SeccompRuleFileIO: "c_cpp_file_io",
	}
	l.configMap["cpp"] = config{
		SrcName:           "main.cpp",
		ExeName:           "main",
		MaxCpuTime:        10000,
		MaxRealTime:       20000,
		MaxMemory:         1024 * 1024 * 1024,
		CompilerPath:      "/usr/bin/g++",
		CompileArgs:       "-DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c++14 {src_path} -lm -o {exe_path}",
		RunCommand:        "{exePath}",
		SeccompRule:       "c_cpp",
		SeccompRuleFileIO: "c_cpp_file_io",
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
	return file.MakeFilePath(dir, conf.SrcName).String(), nil
}

func (l *LanguageConfig) MakeExePath(dir string, language string) (string, error) {
	conf, err := l.Get(language)
	if err != nil {
		return "", fmt.Errorf("failed to make exepath: %w", err)
	}
	return file.MakeFilePath(dir, conf.ExeName).String(), nil
}

func (l *LanguageConfig) MakeCompileArgSlice(srcPath string, exePath string, language string) ([]string, error) {

	conf, err := l.Get(language)
	if err != nil {
		return nil, fmt.Errorf("failed to make argslice: %w", err)
	}
	args := strings.Replace(conf.CompileArgs, "{srcPath}", srcPath, 1)
	args = strings.Replace(args, "{exePath}", exePath, 1)
	return strings.Split(args, " "), nil
}
