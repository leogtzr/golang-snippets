package main

import (
	"fmt"
)

func main() {
	// rand.Seed(time.Now().UnixNano())
	// fileName := "/home/leo/tamal con arroz.txt"
	// fileName = strings.ReplaceAll(path.Base(fileName), " ", "_")

	// finalName := fmt.Sprintf("%d-%s", rand.Intn(100), fileName)
	// fmt.Println(finalName)
	x := true
	fmt.Println(x)

	x = !x
	fmt.Println(x)

	x = !x
	fmt.Println(x)

	x = !x
	fmt.Println(x)
}
