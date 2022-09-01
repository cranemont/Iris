package main

import (
	"bytes"
	"fmt"
	"regexp"
	"unicode"
)

func reg(a []byte) {
	r := regexp.MustCompile("[ \r\t]+(\r\n?|\n)")
	b := r.ReplaceAll(a, []byte("$1"))
	fmt.Println(a, "\n", b)
	fmt.Println(string(b))
}

func tr(r rune) bool {
	return r != '\n' && unicode.IsSpace(r)
}

func te(a []byte) {
	// a := []byte("1 2   \t\r\n  1 2   \n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n")
	b := bytes.Split(a, []byte("\n"))

	for idx, val := range b {
		b[idx] = bytes.TrimRightFunc(val, unicode.IsSpace)
	}
	fmt.Println(a, "\n", bytes.Join(b, []byte("\n")))
	fmt.Println(string(bytes.Join(b, []byte("\n"))))
}

func main() {
	reg([]byte("1 2   \t\r\n  1 2   \n"))
	te([]byte("1 2   \t\r\n  1 2   \n"))
}
