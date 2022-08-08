package dto

type Receive struct {
}

type Send struct {
}

type CompileRequestDto struct {
	Code     string
	Language string
}

type CompileResponseDto struct {
	Result     string // constants에 정의된 result code로
	CompileKey string
}

type RunRequestDto struct {
}
