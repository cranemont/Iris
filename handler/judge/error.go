package judge

// code
const (
	SUCCESS = 0 + iota
	COMPILE_FAILED
	RUN_FAILED
	TESTCASE_GET_FAILED
)

type JudgeError struct {
	mode int
	code int
}

func (j *JudgeError) Error() string {
	return "ERRRRRRRR!@!@"
}
