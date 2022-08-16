package main

import (
	"fmt"

	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/constants/language"
	"github.com/cranemont/judge-manager/event"
	"github.com/cranemont/judge-manager/judger"
	"github.com/cranemont/judge-manager/judger/config"
	"github.com/cranemont/judge-manager/manager"
	"github.com/cranemont/judge-manager/mq"
	"github.com/cranemont/judge-manager/task"
)

func main() {

	eventMap := make(map[string](chan interface{}))
	eventListener := event.NewTaskEventListener(eventMap)
	eventEmitter := event.NewEventEmitter(eventMap)
	eventManager := manager.NewEventManager(eventMap, eventListener, eventEmitter)

	eventManager.Listen(constants.TASK_EXITED, "PublishResult")
	// eventManager.Listen(constants.TASK_EXEC, "Exec")

	sandbox := judger.NewSandbox()

	compileOption := config.CompileOption{}
	runOption := config.RunOption{}

	compiler := judger.NewCompiler(sandbox, &compileOption)
	runner := judger.NewRunner(sandbox, &runOption)
	grader := judger.NewGrader()

	judger := judger.NewJudger(compiler, runner, grader)
	judgeManager := manager.NewJudgeManager(judger, eventEmitter)

	// go task event listener
	// go global error hander

	for {
		var input string
		fmt.Scanln(&input)

		submissionDto := mq.SubmissionDto{
			Code:      "#include <stdio.h>\n\nint main (void) {\nprintf('Hello world!');\nreturn 0;\n}\n",
			Language:  language.C,
			ProblemId: input,
			Limits: mq.Limits{
				Time:   "TIMELIMIT",
				Memory: "MEMORYLIMIT",
			},
		}

		// 아래 과정을 event trigger로 넣어주면?
		// eventManager.Listen(constants.TASK_EXEC, "Exec") 여기서 listen하고 있는 event를
		// eventManager.Emit(constants.TASK_EXEC, submissionDto) 해주고
		// taskManager에서 DTO로 task만든다음에 자기 map에 등록하고(관리용), go judgeManager.Exec(task) 해주거나 아니면 또 trigger해주면?
		// 성능은 좀 떨어져도 일관될것같은데?
		task := task.NewTask(submissionDto)
		// register task to event manager
		// 등록해두고 종료되었음을 감지

		// 큐를 넣을거면 여기서 관리(taskManager에서)
		go judgeManager.Exec(task)
	}
	// 여기서 rabbitMQ consumer가 돌고
	// 메시지 수신시 채점자 호출
}
