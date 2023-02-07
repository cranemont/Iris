package testcase

import (
	"fmt"

	"github.com/cranemont/iris/src/service/cache"
)

type Manager interface {
	GetTestcase(problemId string) (Testcase, error)
	UnMarshal(data []byte) (Testcase, error)
}

type manager struct {
	source Source
	cache  cache.Cache
}

func NewManager(s Source, c cache.Cache) *manager {
	return &manager{source: s, cache: c}
}

func (m *manager) GetTestcase(problemId string) (Testcase, error) {
	isExist, err := m.cache.IsExist(problemId)
	if err != nil {
		return Testcase{}, fmt.Errorf("GetTestcase: %w", err)
	}
	if !isExist {
		testcase, err := m.source.GetTestcase(problemId)
		if err != nil {
			return Testcase{}, fmt.Errorf("get testcase: %w", err)
		}

		err = m.cache.Set(problemId, testcase)
		if err != nil {
			return Testcase{}, fmt.Errorf("cache set: %w", err)
		}
		return testcase, nil
	}
	data, err := m.cache.Get(problemId)
	if err != nil {
		return Testcase{}, fmt.Errorf("testcase: %s: %w", problemId, err)
	}
	testcase, err := m.UnMarshal(data)
	if err != nil {
		return Testcase{}, fmt.Errorf("testcase: %w", err)
	}
	return testcase, nil
}

func (m *manager) UnMarshal(data []byte) (Testcase, error) {
	// validate testcase
	testcase := Testcase{}
	err := testcase.UnmarshalBinary(data)
	if err != nil {
		return Testcase{}, err
	}
	return testcase, nil
}
