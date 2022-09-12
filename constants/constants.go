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
const BASE_DIR_DEV = "/go/src/github.com/cranemont/judge-manager"
const OUTPUT_PATH_DEV = BASE_DIR_DEV + "/output"
const LIBJUDGER_PATH_DEV = BASE_DIR_DEV + "/libjudger.so"
const JAVA_POLICY_PATH_DEV = BASE_DIR_DEV + "/policy/java_policy"

// const BASE_DIR = "/go/src/workspace/results"
const BASE_FILE_MODE = 0755

const BASE_DIR_PROD = "/app/sandbox"
const OUTPUT_PATH_PROD = BASE_DIR_PROD + "/output"
const LIBJUDGER_PATH_PROD = BASE_DIR_PROD + "/libjudger.so"
const JAVA_POLICY_PATH_PROD = BASE_DIR_PROD + "/policy/java_policy"
