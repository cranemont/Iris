package main

import (
	"bytes"
	"regexp"
	"testing"
	"unicode"
)

func BenchmarkMyFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		myFunc([]byte("1 2   \t\r\n  1 2   \n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n"))
	}
}

func BenchmarkRegex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		regex([]byte("1 2   \t\r\n  1 2   \n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n1 2   \t\r\n"))
	}
}

func myFunc(a []byte) []byte {
	b := bytes.Split(a, []byte("\n"))

	for idx, val := range b {
		b[idx] = bytes.TrimRightFunc(val, unicode.IsSpace)
	}
	return bytes.Join(b, []byte("\n"))
}

func regex(a []byte) []byte {
	r := regexp.MustCompile("[ \r\t]+(\r\n?|\n)")
	return r.ReplaceAll(a, []byte("$1"))
}
