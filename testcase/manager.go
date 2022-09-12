package testcase

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cranemont/judge-manager/cache"
)

// FIXME: judge 안에 들어가는게 더 맞을듯
type Manager interface {
	GetTestcase(problemId string) (Testcase, error)
	UnMarshal(data []byte) (Testcase, error)
}

type manager struct {
	cache cache.Cache
	serverUrl string
	token string
}

func NewManager(cache cache.Cache) *manager {
	return &manager{
		cache: cache, 
		serverUrl: "http://localhost:3000/testcase/problem/",//os.Getenv("TESTCASE_SERVER_URL"),
		token: "ro==" ,//os.Getenv("TESTCASE_SERVER_AUTH_TOKEN"),
	}
}

func (m *manager) GetTestcase(problemId string) (Testcase, error) {
	if !m.cache.IsExist(problemId) {
		fmt.Println("Tc does not exist")
		testcase, err := m.GetTestcaseFromServer(problemId)
		if err != nil {
			return Testcase{}, fmt.Errorf("failed to get testcase from server: %w", err)
		}
		m.cache.Set(problemId, testcase)
		return testcase, nil
	}
	data := m.cache.Get(problemId)
	testcase, err := m.UnMarshal(data)
	if err != nil {
		return Testcase{}, fmt.Errorf("testcase byte to slice failed: %w", err)
	}
	return testcase, nil
}

func (m *manager) GetTestcaseFromServer(problemId string) (Testcase, error) {
	// FIXME: timeout 설정
	req, err := http.NewRequest("GET", m.serverUrl+problemId, nil)
    if err != nil {
        return Testcase{}, fmt.Errorf("failed to create request: %w\n", err)
    }
    req.Header.Add("judge-server-token", m.token)
 
    client := &http.Client{}
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
	fmt.Println(bytes)
	testcase, err := m.UnMarshal(bytes)
	if err != nil {
		return Testcase{}, fmt.Errorf("invalid testcase data: %w\n", err)
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
