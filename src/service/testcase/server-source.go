package testcase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cranemont/iris/src/common/constants"
)

type server struct {
	url   string
	token string
}

func NewServer(url, token string) *server {
	return &server{
		url:   url,
		token: token,
	}
}

func (s *server) GetTestcase(problemId string) (Testcase, error) {
	testcase, err := s.getFromServer(problemId)
	if err != nil {
		return Testcase{}, fmt.Errorf("testcase: %w", err)
	}
	return testcase, nil
}

func (s *server) getFromServer(problemId string) (Testcase, error) {
	req, err := http.NewRequest("GET", s.url+problemId, nil)
	if err != nil {
		return Testcase{}, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Add(constants.TOKEN_HEADER, s.token)

	client := &http.Client{Timeout: constants.TESTCASE_GET_TIMEOUT * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Testcase{}, fmt.Errorf("http client error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Testcase{}, fmt.Errorf("status code is not 200: code: %d", resp.StatusCode)
	}

	// 결과 출력
	bytes, _ := ioutil.ReadAll(resp.Body)
	e := []Element{}
	if err := json.Unmarshal(bytes, &e); err != nil {
		return Testcase{}, fmt.Errorf("invalid testcase data: %w", err)
	}

	return Testcase{Data: e}, nil
}
