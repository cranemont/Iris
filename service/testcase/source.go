package testcase

type Source interface {
	GetTestcase(problemId string) (Testcase, error)
}
