package result

type ChResult struct {
	Err   error
	Order int
	Data  interface{}
}
