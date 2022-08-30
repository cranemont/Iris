package testcase

import (
	"fmt"
	"log"

	"github.com/cranemont/judge-manager/cache"
	"github.com/cranemont/judge-manager/common/dto"
)

type TestcaseManager interface {
	// GetTestcase(problemId string) (*Testcase, error)
	GetTestcase(out chan<- dto.GoResult, problemId string)
	CreateTestcaseFromByteSlice(data []byte) (Testcase, error)
}

type testcaseManager struct {
	cache cache.Cache
}

func NewTestcaseManager(cache cache.Cache) *testcaseManager {
	return &testcaseManager{cache}
}

func (t *testcaseManager) GetTestcase(out chan<- dto.GoResult, problemId string) {
	if !t.cache.IsExist(problemId) {
		fmt.Println("Tc does not exist")
		// http get
		// cache set
		// 임시로 생성
		testcase := Testcase{[]TestcaseElement{{In: "1 1\n", Out: "1 1\n"}, {In: "2 2\n", Out: "2 2\n"}}}
		t.cache.Set(problemId, testcase)
		out <- dto.GoResult{Data: testcase}
		return
	}
	data := t.cache.Get(problemId)
	testcase, err := t.CreateTestcaseFromByteSlice(data)
	if err != nil {
		log.Println("Error when getting testcase: ", err)
		out <- dto.GoResult{Err: err}
	}
	out <- dto.GoResult{Data: testcase}
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
