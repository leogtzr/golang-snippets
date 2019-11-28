package main

import (
	"fmt"
	"strings"
)

func removeFirstChar(s string) string {
	if len(s) > 0 && strings.HasPrefix(s, "n") {
		return s[1:]
	}
	return s
}

func main() {
	s := "nParece que simon ... "
	fmt.Println(s)
	fmt.Println(removeFirstChar(s))
}
