package handler

import "encoding/json"

type Code int

// FIXME: egress에 정의되어 있어야 함
type Result struct {
	StatusCode Code        `json:"statusCode"`
	Data       interface{} `json:"data"`
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

type Handler interface {
	Handle(data interface{}) error
}

func DefaultResult() string {
	res, err := json.Marshal(Result{StatusCode: INVALID_MODE})
	if err != nil {
		panic(err)
	}
	return string(res)
}
