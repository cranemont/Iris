package handler

import "encoding/json"

type Code int

type Result struct {
	StatusCode Code        `json:"statusCode"`
	Data       interface{} `json:"data"`
}

type Handler interface {
	Handle(data interface{}) error
}

const (
	SUCCESS Code = 0 + iota
	COMPILE_ERROR
	SANDBOX_ERROR
	TESTCASE_GET_FAILED
	INVALID_TESTCASE
	INVALID_MODE
	INTERNAL_SERVER_ERROR
)

func DefaultResult() string {
	res, err := json.Marshal(Result{StatusCode: INVALID_MODE})
	if err != nil {
		panic(err)
	}
	return string(res)
}
