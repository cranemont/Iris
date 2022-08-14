package main

import (
	"fmt"
	"sync"

	"github.com/cranemont/judge-manager/constants/language"
	"github.com/cranemont/judge-manager/judge"
	"github.com/cranemont/judge-manager/judge/manager"
	"github.com/cranemont/judge-manager/mq"
)

func main() {

	sandbox := judge.NewSandbox()
	compiler := judge.NewCompiler(sandbox)
	runner := judge.NewRunner(sandbox)
	submissionManager := manager.NewManager(compiler, runner)

	// submissionDto := mq.SubmissionDto{
	// 	Code:      "#include <stdio.h>\n\nint main (void) {\nprintf('Hello world!');\nreturn 0;\n}\n",
	// 	Language:  language.C,
	// 	ProblemId: "1",
	// }
	// submissionManager.Judge(submissionDto)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
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
		// 큐는 여기서 관리해줘야지.. 아래는 그냥 일만 하고...
		// 아래 고루틴 Judge에 채널 넘겨줘서 done으로 고루틴 관리? 혹은 group?
		// 아니지 얘는 메소드 호출이라니까? 얘를 고루틴으로 돌리는건 별 문제가 아니다
		wg.Add(1)
		go submissionManager.Judge(submissionDto, &wg)
	}
	// 여기서 rabbitMQ consumer가 돌고
	// 메시지 수신시 채점자 호출
}
