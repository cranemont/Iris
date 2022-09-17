package testcase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/cranemont/judge-manager/cache"
	"github.com/cranemont/judge-manager/constants"
)

type Manager interface {
	GetTestcase(problemId string) (Testcase, error)
	UnMarshal(data []byte) (Testcase, error)
}

type manager struct {
	cache     cache.Cache
	serverUrl string
	token     string
}

func NewManager(cache cache.Cache) *manager {
	return &manager{
		cache:     cache,
		serverUrl: os.Getenv("TESTCASE_SERVER_URL"),
		token:     os.Getenv("TESTCASE_SERVER_AUTH_TOKEN"),
	}
}

func (m *manager) GetTestcase(problemId string) (Testcase, error) {
	isExist, err := m.cache.IsExist(problemId)
	if err != nil {
		return Testcase{}, fmt.Errorf("GetTestcase: %w", err)
	}
	if !isExist {
		testcase, err := m.GetTestcaseFromServer(problemId)
		if err != nil {
			return Testcase{}, fmt.Errorf("failed to get testcase from server: %w", err)
		}
		// temp data
		// testcase := Testcase{
		// 	[]Element{
		// 		{Id: "problem:1:1", In: "1\n", Out: "1\n"},
		// 		{Id: "problem:1:2", In: "22\n", Out: "222\n"},
		// 	},
		// }
		err = m.cache.Set(problemId, testcase)
		if err != nil {
			return Testcase{}, fmt.Errorf("GetTestcase: %w", err)
		}
		return testcase, nil
	}
	data, err := m.cache.Get(problemId)
	if err != nil {
		return Testcase{}, fmt.Errorf("GetTestcase: %s: %w", problemId, err)
	}
	testcase, err := m.UnMarshal(data)
	if err != nil {
		return Testcase{}, fmt.Errorf("testcase byte to slice failed: %w", err)
	}
	return testcase, nil
}

func (m *manager) GetTestcaseFromServer(problemId string) (Testcase, error) {
	req, err := http.NewRequest("GET", m.serverUrl+problemId, nil)
	if err != nil {
		return Testcase{}, fmt.Errorf("failed to create http request: %w\n", err)
	}
	req.Header.Add(constants.TOKEN_HEADER, m.token)

	client := &http.Client{Timeout: constants.TESTCASE_GET_TIMEOUT * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Testcase{}, fmt.Errorf("http client error: %w\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Testcase{}, fmt.Errorf("status code is not 200:\n code: %d\n", resp.StatusCode)
	}

	// 결과 출력
	bytes, _ := ioutil.ReadAll(resp.Body)
	testcaseElements := []Element{}
	if err := json.Unmarshal(bytes, &testcaseElements); err != nil {
		return Testcase{}, fmt.Errorf("invalid testcase data: %w\n", err)
	}

	return Testcase{Data: testcaseElements}, nil
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
