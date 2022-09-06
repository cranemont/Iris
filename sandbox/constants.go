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

const (
	CompileLogPath = "./log/compile/log.out"
	RunLogPath     = "./log/run/log.out"
	CompileOutFile = "compile.out"
)
