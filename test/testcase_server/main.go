package main

import (
	"fmt"
	"net/http"
	"strings"
)

func parseId(path string) string {
	return strings.Split(path, "/")[2]
}

func testcaseHandler(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	fmt.Println(parseId(p))
}

func main() {
	http.HandleFunc("/problem/", testcaseHandler)
	http.ListenAndServe(":30000", nil)
}
