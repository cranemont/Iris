package config

// 언어별 설정 관리..
// 각 DTO로 변환 과정 수행
type CompileOption struct {
	// Get(language) -> constants패키지에서 설정값 가져옴
}

// enum으로 바꾸기
func (c *CompileOption) Get(lang string) string {
	return "test"
}
