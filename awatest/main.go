package main

import (
	"fmt"
	"time"
)

func event(s string) {
	fmt.Println("Event ", s, " waiting...")
	time.Sleep(time.Second * 3)
	fmt.Println("Event ", s, " Done!")
}

type NilSt struct {
	C string
}
type NoField struct {
	A string
	B string
	N *NilSt
}

func main() {
	t := NoField{A: "aaa"}
	if t.N == nil {
		fmt.Println("VLAIAIA")
	}
	for {
		var s string
		fmt.Print("input anything to run")
		fmt.Scanln(&s)
		go event(s)
	}
}
