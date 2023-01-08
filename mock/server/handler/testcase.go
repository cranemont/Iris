package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/cranemont/judge-manager/mock/server/handler/response"
)

type testcase struct {
	r response.Responser
}

func NewTestcaseHandler(r response.Responser) *testcase {
	return &testcase{r}
}

func (t *testcase) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	params, _ := req.Context().Value("params").(map[string]string)

	// FIXME: 경로 환경변수로 설정 가능하도록 수정
	// FIXME: 테스트 가능한 구조로 분리
	path, _ := filepath.Abs("../testcase")
	data, err := os.ReadFile(path + "/" + params["id"] + ".json")
	if err != nil {
		t.r.Error(w, "failed to read testcase file", http.StatusNotFound)
		return
	}

	t.r.Ok(w, response.Data{Message: "ok", Result: data}, http.StatusOK)
}
