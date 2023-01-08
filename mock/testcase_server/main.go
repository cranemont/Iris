package main

import (
	"net/http"

	"github.com/cranemont/judge-manager/mock/testcase_server/handler"
	"github.com/cranemont/judge-manager/mock/testcase_server/middleware"
	"github.com/cranemont/judge-manager/mock/testcase_server/router"
	"github.com/cranemont/judge-manager/mock/testcase_server/router/method"
)

func main() {
	r := router.NewRouter()
	r.Handle(method.GET, "/problem/:id/testcase",
		middleware.Adapt(
			http.HandlerFunc(handler.TestcaseHandler),
			middleware.SetContentType(),
		),
	)
	http.ListenAndServe(":30000", r)
}
