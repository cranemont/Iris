package middleware

import (
	"fmt"
	"net/http"
)

func Response() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Println("before")
			defer fmt.Println("after")
			h.ServeHTTP(w, req)
		})
	}
}
