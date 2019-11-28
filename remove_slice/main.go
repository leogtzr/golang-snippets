package main

import "fmt"

func paginate(x []string, skip int, size int) []string {
	if skip > len(x) {
		skip = len(x)
	}

	end := skip + size
	if end > len(x) {
		end = len(x)
	}

	return x[skip:end]
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

func main() {
	s := make([]string, 0)
	for i := 0; i < 100; i++ {
		s = append(s, fmt.Sprintf("e%d", i))
	}

	// a := paginate(s, 0, 10)
	x := 0
	for i := 0; i < 15; i++ {
		a := paginate(s, x, 10)
		fmt.Println(a)
		x += 10
	}

	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~>")
	fmt.Println(s)

}
