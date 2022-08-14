package mq

type Limits struct {
	Time   string
	Memory string
}

type Testcases struct {
	total  int
	input  []string
	output []string
}

func (t *Testcases) isValid() bool {
	// Input과 output 개수 같은지 확인
	return true
}

type SubmissionDto struct {
	Code      string
	Language  string
	ProblemId string
	Limits    Limits
	Testcases Testcases
}
