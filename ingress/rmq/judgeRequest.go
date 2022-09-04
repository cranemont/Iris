package rmq

type JudgeRequest struct {
	Id          string
	Code        string
	Language    string
	ProblemId   string
	TimeLimit   int
	MemoryLimit int
}
