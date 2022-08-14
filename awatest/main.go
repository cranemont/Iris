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

func main() {
	for {
		var s string
		fmt.Print("input anything to run")
		fmt.Scanln(&s)
		go event(s)
	}
}
