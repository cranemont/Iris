package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/cranemont/judge-manager/cache"
	"github.com/cranemont/judge-manager/constants/language"
	"github.com/cranemont/judge-manager/egress"
	"github.com/cranemont/judge-manager/handler"
	"github.com/cranemont/judge-manager/handler/judge"
	"github.com/cranemont/judge-manager/ingress/rmq"
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

	languageConfig := sandbox.LanguageConfig{}
	languageConfig.Init()

	compiler := sandbox.NewCompiler(&languageConfig)
	runner := sandbox.NewRunner(&languageConfig)

	// eventMap := make(map[string](chan interface{}))
	// eventEmitter := event.NewEventEmitter(eventMap)

	ctx := context.Background()
	cache := cache.NewCache(ctx)
	testcaseManager := testcase.NewTestcaseManager(cache)

	judger := judge.NewJudger(
		compiler,
		runner,
		testcaseManager,
	)
	// specialJudger
	// customTestcaseRunner 만들어서 같이 넣어주기

	publisher := egress.NewRmqPublisher()

	handler := handler.NewHandler(
		judger,
		&languageConfig,
		publisher,
	)

	// judgeEventListener := event.NewListener(eventMap, judgeEventHander)

	// judgeEventManager := event.NewManager(
	// 	eventMap, judgeEventListener, eventEmitter,
	// )

	// judgeEventManager.Listen(constants.TASK_EXEC, "OnExec")
	// judgeEventManager.Listen(constants.TASK_EXITED, "OnExit")

	// go task event listener
	// go global error hander

	for {
		var input string
		fmt.Scanln(&input)

		submissionDto := rmq.SubmissionDto{
			// Code: "#include <stdio.h>\n\nint main (void) {\n  printf(\"Hello world!\\n\");\n  char buf[100];\n  scanf(\"%s\", buf);\n  printf(\"%s\\n\", buf);\n  return 0;\n}\n",
			// Code: "#include <stdio.h>\n\nint main (void) {\nprintf(\"Hello world!\");\nreturn 0;\n}\n",
			// Code:      "#include <stdio.h>\n\nint main (void) {\nwhile(1) {printf(\"Hello world!\");}\nreturn 0;\n}\n",
			Code:        "#include <stdio.h>\n\nint main (void) {\nprintf(\"1 1  \t\\n\");\n\nreturn 0;\n}\n",
			Language:    language.C,
			ProblemId:   input,
			TimeLimit:   1000,
			MemoryLimit: 256 * 1024 * 1024,
		}

		// FIXME: task도 handler에서 만들어야 하나? 그래야 라우팅도 되니까
		// task := judge.NewTask(submissionDto)
		// go judgeEventHander.OnExec(task) //<- 이방법이 더 빠름
		// 라우터를 하나 만들고, SPJ인지, 등등 판단
		go handler.Handle("Judge", submissionDto)
		// judgeEventManager.Dispatch(constants.TASK_EXEC, task) // <- 이벤트를 사용하는 일관된 방법
	}
	// 여기서 rabbitMQ consumer가 돌고
	// 메시지 수신시 채점자 호출
}
