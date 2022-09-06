package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/cranemont/judge-manager/cache"
	"github.com/cranemont/judge-manager/egress"
	"github.com/cranemont/judge-manager/handler"
	"github.com/cranemont/judge-manager/handler/judge"
	"github.com/cranemont/judge-manager/ingress/rmq"
	"github.com/cranemont/judge-manager/router"
	"github.com/cranemont/judge-manager/sandbox"
	"github.com/cranemont/judge-manager/testcase"
)

// func init() {
// 	http.HandleFunc("/debug/pprof/", Index) // Profile Endpoint for Heap, Block, ThreadCreate, Goroutine, Mutex
// 	http.HandleFunc("/debug/pprof/cmdline", Cmdline)
// 	http.HandleFunc("/debug/pprof/profile", Profile) // Profile Endpoint for CPU
// 	http.HandleFunc("/debug/pprof/symbol", Symbol)
// 	http.HandleFunc("/debug/pprof/trace", Trace)
// }

func main() {

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	ctx := context.Background()
	cache := cache.NewCache(ctx)
	testcaseManager := testcase.NewTestcaseManager(cache)

	judger := judge.NewJudger(
		testcaseManager,
	)

	judgeHandler := handler.NewJudgeHandler(judger)
	// specialJudger
	// customTestcaseRunner 만들어서 같이 넣어주기

	publisher := egress.NewRmqPublisher()

	judgeRouter := router.NewRouter(
		judgeHandler,
		publisher,
	)

	for {
		var input string
		fmt.Scanln(&input)

		submissionDto := rmq.JudgeRequest{
			// Code: "#include <stdio.h>\n\nint main (void) {\n  printf(\"Hello world!\\n\");\n  char buf[100];\n  scanf(\"%s\", buf);\n  printf(\"%s\\n\", buf);\n  return 0;\n}\n",
			// Code: "#include <stdio.h>\n\nint main (void) {\nprintf(\"Hello world!\");\nreturn 0;\n}\n",
			// Code: "#include <stdio.h>\n\nint main (void) {\nwhile(1) {printf(\"Hello world!\");}\nreturn 0;\n}\n",
			Code:        "#include <stdio.h>\n\nint main (void) {\nprintf(\"1 1  \t\\n\");\n\nreturn 0;\n}\n",
			Language:    sandbox.C,
			ProblemId:   input,
			TimeLimit:   1000,
			MemoryLimit: 256 * 1024 * 1024,
		}
		go judgeRouter.Route(router.JUDGE, submissionDto)
	}
	// 여기서 rabbitMQ consumer가 돌고
	// 메시지 수신시 채점자 호출
}
