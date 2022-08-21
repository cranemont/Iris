package mq

import "github.com/cranemont/judge-manager/testcase"

// TODO: manager패키지로 옮기기(task 안으로)
type Limit struct {
	Time   string
	Memory string
}

type SubmissionDto struct {
	Id        string
	Code      string
	Language  string
	ProblemId string
	Limit     Limit
	Testcase  testcase.Testcase
}
