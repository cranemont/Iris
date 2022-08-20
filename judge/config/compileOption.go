package config

// 언어별 설정 관리..
// 각 DTO로 변환 과정 수행

type option struct {
	SrcName      string
	ExeName      string
	MaxCpuTime   int
	MaxRealTime  int
	MaxMemory    int
	CompilerPath string
	Args         string
}
type CompileOption struct {
	// Get(language) -> constants패키지에서 설정값 가져옴
	optionMap map[string]*option
}

// enum으로 바꾸기
func (c *CompileOption) Get(lang string) *option {
	return c.optionMap[lang]
}

func (c *CompileOption) Init() {
	c.optionMap = make(map[string]*option)
	// srcPath, exePath는 base dir + task dir
	c.optionMap["c"] = &option{
		SrcName:      "main.c",
		ExeName:      "main",
		MaxCpuTime:   3000,
		MaxRealTime:  10000,
		MaxMemory:    256 * 1024 * 1024,
		CompilerPath: "/usr/bin/gcc",
		Args:         "-DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c11 {srcPath} -lm -o {exePath}",
	}
	c.optionMap["cpp"] = &option{
		SrcName:      "main.cpp",
		ExeName:      "main",
		MaxCpuTime:   10000,
		MaxRealTime:  20000,
		MaxMemory:    1024 * 1024 * 1024,
		CompilerPath: "/usr/bin/g++",
		Args:         "-DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c++14 {src_path} -lm -o {exe_path}",
	}
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
