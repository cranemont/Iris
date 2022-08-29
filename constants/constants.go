package constants

const (
	Success = 1 + iota
	Fail
)

const DIR_NAME_LEN = 16
const MAX_SUBMISSION = 10
const EVENT_CHAN_SIZE = 10
const TASK_EXEC = "Execute"
const TASK_EXITED = "Exited"
const PUBLISH_RESULT = "Publish"

// const BASE_DIR = "./results"
const BASE_DIR = "/go/src/github.com/cranemont/judge-manager/results"

// const BASE_DIR = "/go/src/workspace/results"
const BASE_FILE_MODE = 0755
