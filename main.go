package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/cranemont/judge-manager/cache"
	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/constants/language"
	"github.com/cranemont/judge-manager/event"
	"github.com/cranemont/judge-manager/judge"
	judgeEvent "github.com/cranemont/judge-manager/judge/event"
	"github.com/cranemont/judge-manager/mq"
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

	libjudger := sandbox.NewSandbox()

	languageConfig := sandbox.LanguageConfig{}
	languageConfig.Init()
	// runOption := config.RunOption{}

	compiler := sandbox.NewCompiler(libjudger, &languageConfig)
	runner := sandbox.NewRunner(libjudger, &languageConfig)
	grader := judge.NewGrader()

	eventMap := make(map[string](chan interface{}))
	eventEmitter := event.NewEventEmitter(eventMap)

	ctx := context.Background()
	cache := cache.NewCache(ctx)
	testcaseManager := testcase.NewTestcaseManager(cache)

	judger := judge.NewJudger(
		compiler,
		runner,
		grader,
		testcaseManager,
	)

	judgeEventHander := judgeEvent.NewHandler(
		judger,
		eventEmitter,
		&languageConfig,
	)

	judgeEventListener := event.NewListener(eventMap, judgeEventHander)

	judgeEventManager := event.NewManager(
		eventMap, judgeEventListener, eventEmitter,
	)

	judgeEventManager.Listen(constants.TASK_EXEC, "OnExec")
	judgeEventManager.Listen(constants.TASK_EXITED, "OnExit")

	// go task event listener
	// go global error hander

	for {
		var input string
		fmt.Scanln(&input)

		submissionDto := mq.SubmissionDto{
			// Code: "#include <stdio.h>\n\nint main (void) {\n  printf(\"Hello world!\\n\");\n  char buf[100];\n  scanf(\"%s\", buf);\n  printf(\"%s\\n\", buf);\n  return 0;\n}\n",
			// Code: "#include <stdio.h>\n\nint main (void) {\nprintf(\"Hello world!\");\nreturn 0;\n}\n",
			// Code:      "#include <stdio.h>\n\nint main (void) {\nwhile(1) {printf(\"Hello world!\");}\nreturn 0;\n}\n",
			Code:      "#include <stdio.h>\n\nint main (void) {\nprintf(\"1 1  \t\\n\");\nreturn 0;\n}\n",
			Language:  language.C,
			ProblemId: input,
			Limit: mq.Limit{
				Time:   1000,
				Memory: 256 * 1024 * 1024,
			},
		}

		task := judge.NewTask(submissionDto)
		// go judgeEventHander.OnExec(task) //<- 이방법이 더 빠름
		judgeEventManager.Dispatch(constants.TASK_EXEC, task) // <- 이벤트를 사용하는 일관된 방법
	}
	// 여기서 rabbitMQ consumer가 돌고
	// 메시지 수신시 채점자 호출
}
