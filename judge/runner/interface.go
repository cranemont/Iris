package runner

type Interface interface {
	Run(ch chan bool)
	Result()
}
