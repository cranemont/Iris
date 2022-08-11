package sandbox

import "fmt"

type Sandbox struct {
}

func NewSandbox() *Sandbox {
	return &Sandbox{}
}

func (s Sandbox) Execute(ch chan bool) {
	fmt.Println("Sandbox: func execute")
	ch <- true
}
