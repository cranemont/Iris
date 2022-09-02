package sandbox

const ( // system error, logging
	SUCCESS = 0 - iota
	INVALID_CONFIG
	FORK_FAILED
	PTHREAD_FAILED
	WAIT_FAILED
	ROOT_REQUIRED
	LOAD_SECCOMP_FAILED
	SETRLIMIT_FAILED
	DUP2_FAILED
	SETUID_FAILED
	EXECVE_FAILED
	SPJ_ERROR
)

const ( // testcase 별 result
	WRONG_ANSWER            = -1
	CPU_TIME_LIMIT_EXCEEDED = 0 + iota
	REAL_TIME_LIMIT_EXCEEDED
	MEMORY_LIMIT_EXCEEDED
	RUNTIME_ERROR
	SYSTEM_ERROR
)

var format = formatString{
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
