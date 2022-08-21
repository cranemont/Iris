package main

import (
	"fmt"

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

func main() {

	sandbox := judge.NewSandbox()

	compileOption := config.CompileOption{}
	runOption := config.RunOption{}

	compiler := judge.NewCompiler(sandbox, &compileOption)
	runner := judge.NewRunner(sandbox, &runOption)
	grader := judge.NewGrader()

	judger := judge.NewJudger(compiler, runner, grader)

	eventMap := make(map[string](chan interface{}))
	eventEmitter := event.NewEventEmitter(eventMap)

	cache := cache.NewCache()
	testcaseManager := testcase.NewTestcaseManager(cache)

	judgeController := judge.NewJudgeController(judger, eventEmitter, testcaseManager)
	fileManager := fileManager.NewFileManager()

	judgeEventHander := judgeEvent.NewJudgeEventHandler(judgeController, fileManager, eventEmitter)
	judgeEventHander.RegisterFn()

	judgeEventListener := event.NewEventListener(eventMap, judgeEventHander)

	judgeEventManager := event.NewEventManager(eventMap, judgeEventListener, eventEmitter)

	judgeEventManager.Listen(constants.TASK_EXEC, "OnExec")
	judgeEventManager.Listen(constants.TASK_EXITED, "OnExit")

	// go task event listener
	// go global error hander

	for {
		var input string
		fmt.Scanln(&input)

		submissionDto := mq.SubmissionDto{
			Code:      "#include <stdio.h>\n\nint main (void) {\nprintf(\"Hello world!\");\nreturn 0;\n}\n",
			Language:  language.C,
			ProblemId: input,
			Limit: mq.Limit{
				Time:   "TIMELIMIT",
				Memory: "MEMORYLIMIT",
			},
		}

		task := judge.NewTask(submissionDto)
		// go judgeEventHander.OnExec(task) //<- 이방법이 더 빠름
		go judgeEventManager.Dispatch(constants.TASK_EXEC, task) // <- 이벤트를 사용하는 일관된 방법
	}
	// 여기서 rabbitMQ consumer가 돌고
	// 메시지 수신시 채점자 호출
}
