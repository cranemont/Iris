package judger

import "github.com/cranemont/judge-manager/mq"

type Receive struct {
}

type Send struct {
}

// TODO: Getter, Setter 만들어서 불변 속성으로 관리

type CompileRequestDto struct {
	Code     string
	Language string
}

type CompileResponseDto struct {
	Result     string // constants에 정의된 result code로
	CompileKey string
}

type RunRequestDto struct{}

type RunResponseDto struct{}

// NewJudgeRequestDto 만들어서 데이터 수정 불가능하게 만들기
type JudgeRequestDto struct {
	RunRequestDto *RunRequestDto
	Testcases     *mq.Testcases
	// RunRequestDto, testcase 객체 저장
}

type JudgeResponseDto struct{}
