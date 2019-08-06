package main

import (
	"fmt"
	"strings"
)

const WrapMax = 80

var line1 = "123456789 123456789 123456789 123456789 123456789"

//var line1 = "leonardoses gutierrezes ramirezastrolopitecusramirez"

func needsWrap(line string) bool {
	len := len(line)
	return len+1 > (WrapMax / 2)
}

func countWithoutWhitespaces(words []string) int {
	count := 0
	for _, w := range words {
		count += len(w)
	}
	return count
}

func wrap(line string) string {
	if !needsWrap(line) {
		return line
	}
	fields := strings.Fields(line)
	numberOfWords := len(fields)
	// wrapBetweenWords := len(line)
	countWithoutSpaces := countWithoutWhitespaces(fields)
	wrapLength := WrapMax - countWithoutSpaces
	//fmt.Println(wrapLength)
	// fmt.Printf("countWithoutSpaces = %d\n", countWithoutSpaces)
	// fmt.Printf("wrapLength = %d\n", wrapLength)

	// fmt.Printf("len(line) = %d\n", len(line))
	// fmt.Printf("wrapLength/numberOfWords = %d\n", wrapLength/numberOfWords)
	return fmt.Sprintf("[%s]", strings.Join(fields, strings.Repeat(" ", wrapLength/(numberOfWords-1))))
}

func main() {
	// fmt.Printf("[%d]\n", len(line1))
	//fmt.Println(needsWrap(line1))
	// fmt.Println(wrap(line1))
	fmt.Println(wrap(line1))
}
