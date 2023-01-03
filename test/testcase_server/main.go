package main

import (
	"net/http"

	"github.com/cranemont/judge-manager/test/testcase_server/handler"
	"github.com/cranemont/judge-manager/test/testcase_server/middleware"
	"github.com/cranemont/judge-manager/test/testcase_server/router"
	"github.com/cranemont/judge-manager/test/testcase_server/router/method"
)

func main() {
	r := router.NewRouter()
	r.Handle(method.GET, "/problem/:id/testcase",
		middleware.Adapt(
			http.HandlerFunc(handler.TestcaseHandler),
			middleware.Example(1),
			middleware.Example(2),
			middleware.Example(3),
		),
	)
	http.ListenAndServe(":30000", r)
}
