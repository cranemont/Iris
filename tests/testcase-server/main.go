package main

import (
	"net/http"

	"github.com/cranemont/judge-manager/tests/testcase-server/handler"
	"github.com/cranemont/judge-manager/tests/testcase-server/handler/response"
	"github.com/cranemont/judge-manager/tests/testcase-server/middleware"
	"github.com/cranemont/judge-manager/tests/testcase-server/router"
	"github.com/cranemont/judge-manager/tests/testcase-server/router/method"
)

func main() {
	r := router.NewRouter()
	responser := response.NewResponser()
	testcaseHandler := handler.NewTestcaseHandler(responser)

	r.Handle(method.GET, "/problem/:id/testcase",
		middleware.Adapt(
			testcaseHandler,
			middleware.SetContentType(),
		),
	)

	http.ListenAndServe(":30000", r)
}
