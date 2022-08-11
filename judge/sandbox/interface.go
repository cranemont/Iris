package sandbox

type Interface interface {
	Execute(ch chan bool)
}
