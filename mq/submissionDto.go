package mq

// TODO: manager패키지로 옮기기(task 안으로)
type Limits struct {
	Time   string
	Memory string
}

type Testcases struct {
	total  int
	input  []string
	output []string
}

func (t *Testcases) IsValid() bool {
	// Input과 output 개수 같은지 확인
	return true
}

func (t *Testcases) GetTotal() int {
	return t.total
}

type SubmissionDto struct {
	Id        string
	Code      string
	Language  string
	ProblemId string
	Limits    Limits
	Testcases Testcases
}
