package sandbox

import (
	"fmt"
	"strings"

	"github.com/cranemont/judge-manager/file"
)

// 언어별 설정 관리..
// 각 DTO로 변환 과정 수행

type config struct {
	SrcName               string
	ExeName               string
	MaxCpuTime            int
	MaxRealTime           int
	MaxMemory             int
	CompilerPath          string
	CompileArgs           string
	RunCommand            string
	RunArgs               string
	SeccompRule           string
	SeccompRuleFileIO     string
	MemoeryLimitCheckOnly bool
	env                   []string
}
type LanguageConfig struct {
	configMap map[string]config
}

func (l *LanguageConfig) Init() {
	l.configMap = make(map[string]config)
	// srcPath, exePath는 base dir + task dir
	defaultEnv := []string{"LANG=en_US.UTF-8", "LANGUAGE=en_US:en", "LC_ALL=en_US.UTF-8"}
	l.configMap["c"] = config{
		SrcName:      "main.c",
		ExeName:      "main",
		MaxCpuTime:   3000,              // compile
		MaxRealTime:  10000,             // compile
		MaxMemory:    256 * 1024 * 1024, // compile
		CompilerPath: "/usr/bin/gcc",
		CompileArgs: "-DONLINE_JUDGE " +
			"-O2 -w -fmax-errors=3 -std=c11 " +
			"{srcPath} -lm -o {exePath}",
		RunCommand:            "{exePath}",
		SeccompRule:           "c_cpp",
		SeccompRuleFileIO:     "c_cpp_file_io",
		MemoeryLimitCheckOnly: false,
		env:                   defaultEnv,
	}
	l.configMap["cpp"] = config{
		SrcName:      "main.cpp",
		ExeName:      "main",
		MaxCpuTime:   10000,
		MaxRealTime:  20000,
		MaxMemory:    1024 * 1024 * 1024,
		CompilerPath: "/usr/bin/g++",
		CompileArgs: "-DONLINE_JUDGE " +
			"-O2 -w -fmax-errors=3 " +
			"-std=c++14 {src_path} -lm -o {exe_path}",
		RunCommand:            "{exePath}",
		SeccompRule:           "c_cpp",
		SeccompRuleFileIO:     "c_cpp_file_io",
		MemoeryLimitCheckOnly: false,
		env:                   defaultEnv,
	}
	l.configMap["java"] = config{
		SrcName:      "Main.java",
		ExeName:      "Main",
		MaxCpuTime:   5000,
		MaxRealTime:  10000,
		MaxMemory:    -1,
		CompilerPath: "/usr/bin/javac",
		CompileArgs:  "{src_path} -d {exe_dir} -encoding UTF8",
		RunCommand:   "/usr/bin/java",
		RunArgs: "-cp {exe_dir} " +
			"-XX:MaxRAM={max_memory}k " +
			"-Djava.security.manager " +
			"-Dfile.encoding=UTF-8 " +
			"-Djava.security.policy==/etc/java_policy " +
			"-Djava.awt.headless=true " +
			"Main",
		SeccompRule:           "None",
		MemoeryLimitCheckOnly: true,
		env:                   defaultEnv,
	}
	l.configMap["py2"] = config{
		SrcName:      "solution.py",
		ExeName:      "solution.pyc",
		MaxCpuTime:   3000,
		MaxRealTime:  10000,
		MaxMemory:    1024 * 1024 * 1024,
		CompilerPath: "/usr/bin/python",
		CompileArgs:  "-m py_compile {src_path}",
		RunCommand:   "/usr/bin/python",
		RunArgs:      "{exe_path}",
		SeccompRule:  "general",
		env:          defaultEnv,
	}
	l.configMap["py3"] = config{
		SrcName:      "solution.py",
		ExeName:      "__pycache__/solution.cpython-38.pyc", // TODO: 파이썬 버전 확인
		MaxCpuTime:   3000,
		MaxRealTime:  10000,
		MaxMemory:    128 * 1024 * 1024,
		CompilerPath: "/usr/bin/python",
		CompileArgs:  "-m py_compile {src_path}",
		RunCommand:   "/usr/bin/python3",
		RunArgs:      "{exe_path}",
		SeccompRule:  "general",
		env:          append(defaultEnv, "PYTHONIOENCODING=utf-8"),
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
		return "", fmt.Errorf("failed to get config: %w", err)
	}
	return file.MakeFilePath(dir, conf.SrcName).String(), nil
}

func (l *LanguageConfig) MakeExePath(dir string, language string) (string, error) {
	conf, err := l.Get(language)
	if err != nil {
		return "", fmt.Errorf("failedto get config: %w", err)
	}
	return file.MakeFilePath(dir, conf.ExeName).String(), nil
}

func (l *LanguageConfig) MakeCompileArgSlice(srcPath string, exePath string, language string) ([]string, error) {

	conf, err := l.Get(language)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	args := strings.Replace(conf.CompileArgs, "{srcPath}", srcPath, 1)
	args = strings.Replace(args, "{exePath}", exePath, 1)
	return strings.Split(args, " "), nil
}
