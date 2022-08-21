package testcase

import "github.com/cranemont/judge-manager/cache"

type TestcaseManager interface {
	GetTestcase(problemId string) *Testcase
	CreateTestcase(data interface{}) *Testcase
}

type testcaseManager struct {
	cache cache.Cache
}

func NewTestcaseManager(cache cache.Cache) *testcaseManager {
	return &testcaseManager{cache}
}

func (t *testcaseManager) GetTestcase(problemId string) *Testcase {
	data := t.cache.Get()
	// FIXME: return type interface로 바꾸고 nil로 변경
	if data == "" {
		// get from http
		// set cache
		// set data
	}
	return &Testcase{}
	// return t.CreateTestcase() ????
}

func (t *testcaseManager) CreateTestcase(data interface{}) *Testcase {
	// validate testcase
	// input, output 개수와 total 일치하는지 등
	return &Testcase{}
}
