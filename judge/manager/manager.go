package manager

import (
	"fmt"

	"github.com/cranemont/judge-manager/judge"
	"github.com/cranemont/judge-manager/mq"
)

type Manager struct {
	compiler judge.Compiler
	runner   judge.Runner
}

func NewManager(
	compiler judge.Compiler,
	runner judge.Runner,
) *Manager {
	return &Manager{compiler, runner}
}

func (m *Manager) Judge(submissionDto mq.SubmissionDto) {
	task := NewTask(submissionDto)
	m.compiler.Compile(task.ToCompileRequestDto())

	// testcase managing?
	m.judge(task) // 정말 문제가 없을까?
	fmt.Println("done")
}

func (m *Manager) judge(t *Task) {

	tcNum := t.GetTestcaseNum()
	// ch := make(chan string)
	for i := 0; i < tcNum; i++ {
		// input := fmt.Sprintf("<test input %d>", i)
		// go t.RequestRun(input, ch)
		go m.runner.Run(t.ToRunRequestDto())
	}
	// for i := 0; i < tcNum; i++ {
	// 	result := <-ch
	// 	fmt.Printf("%s Done!\n", result)
	// 	// 여기서 이제 grade 고루틴으로 정리
	// }
	// 여기서도 wait?
}
