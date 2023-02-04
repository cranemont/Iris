package testcase

type preset struct {
}

func NewPreset() *preset {
	return &preset{}
}

func (p *preset) GetTestcase(problemId string) (Testcase, error) {
	return Testcase{
		[]Element{
			{Id: "problem:1:1", In: "1\n", Out: "1\n"},
			{Id: "problem:1:2", In: "22\n", Out: "222\n"},
		},
	}, nil
}
