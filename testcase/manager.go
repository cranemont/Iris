package testcase

import (
	"fmt"

	"github.com/cranemont/judge-manager/cache"
)

// TODO: event의 manager랑 이름 헷갈림, 이름 더 명확하게 바꾸기
type TestcaseManager interface {
	// GetTestcase(problemId string) (*Testcase, error)
	GetTestcase(problemId string) (Testcase, error)
	CreateTestcaseFromByteSlice(data []byte) (Testcase, error)
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
		// cache set
		// 임시로 생성
		testcase := Testcase{[]TestcaseElement{{In: "1 1\n", Out: "1 1\n"}, {In: "2 2\n", Out: "2 2\n"}}}
		t.cache.Set(problemId, testcase)
		return testcase, nil
	}
	data := t.cache.Get(problemId)
	testcase, err := t.CreateTestcaseFromByteSlice(data)
	if err != nil {
		return Testcase{}, fmt.Errorf("testcase byte to slice failed: %w", err)
	}
	return testcase, nil
}

func (t *testcaseManager) CreateTestcaseFromByteSlice(data []byte) (Testcase, error) {
	// validate testcase
	testcase := Testcase{}
	err := testcase.UnmarshalBinary(data)
	if err != nil {
		return Testcase{}, err
	}
	return testcase, nil
}
