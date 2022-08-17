package mq

// TODO: manager패키지로 옮기기(task 안으로)
type Limit struct {
	Time   string
	Memory string
}

type Testcase struct {
	total  int
	input  []string
	output []string
}

func (t *Testcase) IsValid() bool {
	// Input과 output 개수 같은지 확인
	return true
}

func (t *Testcase) GetTotal() int {
	return t.total
}

type SubmissionDto struct {
	Id        string
	Code      string
	Language  string
	ProblemId string
	Limit     Limit
	Testcase  Testcase
}
