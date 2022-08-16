package judger

import "fmt"

type Judger interface {
	Compile(dto *CompileRequestDto)
	Judge(dto *JudgeRequestDto) // run and grade
}

type judger struct {
	compiler Compiler
	runner   Runner
	grader   Grader
}

func NewJudger(compiler Compiler, runner Runner, grader Grader) *judger {
	return &judger{
		compiler: compiler,
		runner:   runner,
		grader:   grader,
	}
}

// err 처리
func (j *judger) Compile(dto *CompileRequestDto) {
	j.compiler.Compile(dto)
}

// err 처리
func (j *judger) Judge(dto *JudgeRequestDto) {
	// run and grade
	tcNum := dto.Testcases.GetTotal()
	ch := make(chan string, tcNum)
	for i := 0; i < tcNum; i++ {
		go j.runner.Run(dto.RunRequestDto)
	}
	for i := 0; i < tcNum; i++ {
		result := <-ch
		fmt.Printf("%s Done!\n", result)
		// 여기서 이제 grade 고루틴으로 정리
	}
	close(ch)
	// close chan
}
