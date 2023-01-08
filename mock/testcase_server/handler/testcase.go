package handler

import (
	"fmt"
	"net/http"
)

func TestcaseHandler(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	fmt.Println(p)
	params, ok := req.Context().Value("params").(map[string]string)
	if ok {
		for k, v := range params {
			fmt.Println(k, v)
		}
	}
	r := Response{Message: "mssage", Data: "data"}
	r.Encode(w)
}
