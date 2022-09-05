package testcase

import (
	"fmt"

	"github.com/cranemont/judge-manager/cache"
)

// FIXME: judge 안에 들어가는게 더 맞을듯
type TestcaseManager interface {
	GetTestcase(problemId string) (Testcase, error)
	UnMarshal(data []byte) (Testcase, error)
}

type testcaseManager struct {
	cache cache.Cache
}

func NewTestcaseManager(cache cache.Cache) *testcaseManager {
	return &testcaseManager{cache}
}

func (t *testcaseManager) GetTestcase(problemId string) (Testcase, error) {
	if !t.cache.IsExist(problemId) {
		fmt.Println("Tc does not exist")
		// http get
		// 임시로 생성
		testcase := Testcase{[]Element{{In: "1 1\n", Out: "1 1\n"}, {In: "2 2\n", Out: "2 2\n"}}}
		t.cache.Set(problemId, testcase)
		return testcase, nil
	}
	data := t.cache.Get(problemId)
	testcase, err := t.UnMarshal(data)
	if err != nil {
		return Testcase{}, fmt.Errorf("testcase byte to slice failed: %w", err)
	}
	return testcase, nil
}

func (t *testcaseManager) UnMarshal(data []byte) (Testcase, error) {
	// validate testcase
	testcase := Testcase{}
	err := testcase.UnmarshalBinary(data)
	if err != nil {
		return Testcase{}, err
	}
	return testcase, nil
}
