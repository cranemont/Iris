package testcase

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
	// return t.total
	return 3
}
