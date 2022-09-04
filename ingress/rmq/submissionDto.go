package rmq

type SubmissionDto struct {
	Id          string
	Code        string
	Language    string
	ProblemId   string
	TimeLimit   int
	MemoryLimit int
}
