package main

import (
	"fmt"
	"net/http"

	"github.com/cranemont/judge-manager/test/testcase_server/router"
	"github.com/cranemont/judge-manager/test/testcase_server/router/method"
)

func testcaseHandler(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	fmt.Println(p)
	params, ok := req.Context().Value("params").(map[string]string)
	if ok {
		for k, v := range params {
			fmt.Println(k, v)
		}
	}
	w.Write([]byte{'d', 'e', 'f'})
}

func main() {
	r := router.NewRouter()
	r.HandleFunc(method.GET, "/problem/:id/testcase", testcaseHandler)
	http.ListenAndServe(":30000", r)
}
