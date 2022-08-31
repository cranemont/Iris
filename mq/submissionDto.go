package mq

import "github.com/cranemont/judge-manager/testcase"

// TODO: manager패키지로 옮기기(task 안으로)
type Limit struct {
	Time   int
	Memory int
}

type SubmissionDto struct {
	Id        string
	Code      string
	Language  string
	ProblemId string
	Limit     Limit
	Testcase  testcase.Testcase // 없을수도 있음. 다른 종류의 struct
}
