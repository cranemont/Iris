package main

import (
	"fmt"
	"sync"

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

	var wg sync.WaitGroup
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

		task := task.NewTask(submissionDto)
		// register task to event manager
		// 등록해두고 종료되었음을 감지

		// 큐를 넣을거면 여기서 관리
		wg.Add(1) // 필요할까? 위에서 register하고 거기서 관리하면?
		go judgeManager.Exec(task, &wg)
	}
	// 여기서 rabbitMQ consumer가 돌고
	// 메시지 수신시 채점자 호출
}
