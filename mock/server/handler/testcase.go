package handler

import (
	"net/http"

	"github.com/cranemont/judge-manager/mock/server/handler/response"
)

type testcase struct {
	r response.Responser
}

func NewTestcaseHandler(r response.Responser) *testcase {
	return &testcase{r}
}

func (t *testcase) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// params, _ := req.Context().Value("params").(map[string]string)
	t.r.Ok(w, response.Data{Message: "ok", Result: "dfsef"}, http.StatusOK)
}
