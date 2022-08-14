package mq

type Limits struct {
	Time   string
	Memory string
}

type SubmissionDto struct {
	Code      string
	Language  string
	ProblemId string
	Limits    Limits
}
