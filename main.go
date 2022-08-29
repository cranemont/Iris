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
	"github.com/cranemont/judge-manager/fileManager"
	"github.com/cranemont/judge-manager/judge"
	"github.com/cranemont/judge-manager/judge/config"
	judgeEvent "github.com/cranemont/judge-manager/judge/event"
	"github.com/cranemont/judge-manager/mq"
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

	sandbox := judge.NewSandbox()

	compileOption := config.CompileOption{}
	// runOption := config.RunOption{}

	compiler := judge.NewCompiler(sandbox, &compileOption)
	runner := judge.NewRunner(sandbox, &compileOption)
	grader := judge.NewGrader()

	eventMap := make(map[string](chan interface{}))
	eventEmitter := event.NewEventEmitter(eventMap)

	ctx := context.Background()
	cache := cache.NewCache(ctx)
	testcaseManager := testcase.NewTestcaseManager(cache)
	fileManager := fileManager.NewFileManager()

	judgeService := judge.NewJudgeService(
		compiler,
		runner,
		grader,
		fileManager,
		testcaseManager,
	)

	judgeEventHander := judgeEvent.NewJudgeEventHandler(judgeService, eventEmitter)
	judgeEventHander.RegisterFn()

	judgeEventListener := event.NewEventListener(eventMap, judgeEventHander)

	judgeEventManager := event.NewEventManager(
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
			Code: "#include <stdio.h>\n\nint main (void) {\nprintf(\"Hello world!\");\nreturn 0;\n}\n",
			// Code:      "#include <stdio.h>\n\nint main (void) {\nwhile(1) {printf(\"Hello world!\");}\nreturn 0;\n}\n",
			Language:  language.C,
			ProblemId: input,
			Limit: mq.Limit{
				Time:   "TIMELIMIT",
				Memory: "MEMORYLIMIT",
			},
		}

		task := judge.NewTask(submissionDto)
		// go judgeEventHander.OnExec(task) //<- 이방법이 더 빠름
		judgeEventManager.Dispatch(constants.TASK_EXEC, task) // <- 이벤트를 사용하는 일관된 방법
	}
	// 여기서 rabbitMQ consumer가 돌고
	// 메시지 수신시 채점자 호출
}
