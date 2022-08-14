package manager

import (
	"fmt"
	"sync"

	"github.com/cranemont/judge-manager/judge"
	"github.com/cranemont/judge-manager/mq"
)

type Manager struct {
	compiler judge.Compiler
	runner   judge.Runner
	// grader
	// config
	// tcManager
	// dirManager
	// errorHander
	// outputGenerator
}

func NewManager(
	// compiler, runner, grader는 하나의 과정이므로 하나로 묶일 수 있음. 그리고 그 객체 안에 config와 dirManager도 포함하겠지
	compiler judge.Compiler,
	runner judge.Runner,
	// grader
	// config
	// dirManager

	// tcManager

	// errorHander
	// outputGenerator -> mq나 다른 외부로 전송
) *Manager {
	return &Manager{compiler, runner}
}

func (m *Manager) Judge(submissionDto mq.SubmissionDto, wg *sync.WaitGroup) {
	defer func() {
		fmt.Print("clean up directory")
		wg.Done()
	}()
	// fmt.Printf("Manager Addr: %p", m)
	// fmt.Printf("Compiler Addr: %p", &m.compiler)
	// fmt.Printf("Runner Addr: %p", &m.runner)

	//
	// task 만들고
	// 디렉토리 만들고
	// 컴파일과 동시에 테스트케이스 가져오기(메모리에 올리기), 동시에 config에서 언어 설정 가져오기... 그것들을 task에 저장하기
	// task의 testcase가 있으면 isValid 체크한다음에 그거 쓰고, 없으면 가져와서 task의 testcase에 저장
	// 이후 m.judge 호출

	task := NewTask(submissionDto)
	// dir 만들기
	m.compiler.Compile(task.ToCompileRequestDto())

	// testcase managing?
	m.judge(task)
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

	// 고민할거는 여기서 바로 output에 쏴줄거냐, 아니면 다른데서 쏴줄거냐?
}
