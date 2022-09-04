package sandbox

const ( // ErrorCode
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

const ( // ResultCode
	RUN_SUCCESS = 0 + iota // this only means the process exited normally
	CPU_TIME_LIMIT_EXCEEDED
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
