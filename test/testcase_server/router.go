package main

import (
	"net/http"
	"strings"
)

func match(pattern, path string) (bool, map[string]string) {
	if pattern == path {
		return true, nil
	}

	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")
	params := make(map[string]string)

	for idx, p := range patterns {
		switch {
		case p == paths[idx]:
		case len(p) > 1 && p[0] == ':':
			params[p[1:]] = paths[idx]
		default:
			return false, nil
		}
	}
	return true, params
}

type router struct {
	handler map[string]http.HandlerFunc
}

func (r *router) GET(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	// pattern 처리, getHandler맵에 handler 등록
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// request 파싱, handler 실행

	http.NotFound(w, req)
}

func (r *router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {

}

// type serverHandler struct {
// 	srv *Server
// }

// func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
// 	handler := sh.srv.Handler
// 	if handler == nil {
// 		handler = DefaultServeMux
// 	}
// 	if req.RequestURI == "*" && req.Method == "OPTIONS" {
// 		handler = globalOptionsHandler{}
// 	}

// 	if req.URL != nil && strings.Contains(req.URL.RawQuery, ";") {
// 		var allowQuerySemicolonsInUse int32
// 		req = req.WithContext(context.WithValue(req.Context(), silenceSemWarnContextKey, func() {
// 			atomic.StoreInt32(&allowQuerySemicolonsInUse, 1)
// 		}))
// 		defer func() {
// 			if atomic.LoadInt32(&allowQuerySemicolonsInUse) == 0 {
// 				sh.srv.logf("http: URL query contains semicolon, which is no longer a supported separator; parts of the query may be stripped when parsed; see golang.org/issue/25192")
// 			}
// 		}()
// 	}

// 	handler.ServeHTTP(rw, req)
// }
