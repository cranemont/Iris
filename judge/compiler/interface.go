package compiler

type Interface interface {
	Compile(ch chan bool)
}
