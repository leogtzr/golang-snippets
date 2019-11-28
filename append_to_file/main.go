package main

import (
	"log"
	"os"
)

const (
	text1 = `
Leo
	Gutiérrez

Ramírez
""
`

	text2 = `
Brenda
	Liliana

Gutiérrez
""
`
)

func appendToFile(content, fileName string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content + "\n")
	return err
}

func main() {
	if err := appendToFile(text1, "text.log"); err != nil {
		log.Fatal(err)
	}

	if err := appendToFile(text2, "text.log"); err != nil {
		log.Fatal(err)
	}

}
