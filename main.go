package main

import (
	"fmt"

	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/constants/language"
	"github.com/cranemont/judge-manager/event"
	"github.com/cranemont/judge-manager/fileManager"
	"github.com/cranemont/judge-manager/judge"
	"github.com/cranemont/judge-manager/judge/config"
	"github.com/cranemont/judge-manager/judge/judgeEvent"
	"github.com/cranemont/judge-manager/mq"
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
	judgeController := judge.NewJudgeController(judger, eventEmitter)
	fileManager := fileManager.NewFileManager()

	judgeEventHander := judgeEvent.NewJudgeEventHandler(judgeController, fileManager, eventEmitter)
	judgeEventListener := judgeEvent.NewJudgeEventListener(eventMap, judgeEventHander)

	judgeEventManager := event.NewEventManager(eventMap, judgeEventListener, eventEmitter)

	judgeEventManager.Listen(constants.TASK_EXEC, "OnExec")
	judgeEventManager.Listen(constants.TASK_EXITED, "OnExit")

	// go task event listener
	// go global error hander

	for {
		var input string
		fmt.Scanln(&input)

		submissionDto := mq.SubmissionDto{
			Code:      "#include <stdio.h>\n\nint main (void) {\nprintf('Hello world!');\nreturn 0;\n}\n",
			Language:  language.C,
			ProblemId: input,
			Limit: mq.Limit{
				Time:   "TIMELIMIT",
				Memory: "MEMORYLIMIT",
			},
		}

		// 아래 과정을 event trigger로 넣어주면?
		// eventManager.Listen(constants.TASK_EXEC, "Exec") 여기서 listen하고 있는 event를
		// eventManager.Emit(constants.TASK_EXEC, submissionDto) 해주고
		// taskManager에서 DTO로 task만든다음에 자기 map에 등록하고(관리용), go judgeManager.Exec(task) 해주거나 아니면 또 trigger해주면?
		// 성능은 좀 떨어져도 일관될것같은데?
		// task := task.NewTask(submissionDto)
		// register task to event manager
		// 등록해두고 종료되었음을 감지

		task := judge.NewTask(submissionDto)
		// go judgeEventHander.OnExec(task) //<- 이방법이 더 빠름
		go judgeEventManager.Dispatch(constants.TASK_EXEC, task) // <- 이벤트를 사용하는 일관된 방법
	}
	// 여기서 rabbitMQ consumer가 돌고
	// 메시지 수신시 채점자 호출
}
