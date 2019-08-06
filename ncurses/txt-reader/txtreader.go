package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/marcusolsson/tui-go"
)

type post struct {
	// username string
	// message  string
	line string
}

// var posts = []post{
// 	{username: "john", message: "hi, what's up?", time: "14:41"},
// 	{username: "jane", message: "not much", time: "14:43"},
// }

var textTestInput = "123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_"

var _posts = []post{
	// {line: "jane"},
}

// ADVANCE ...
const ADVANCE int = 30

var from = 0
var to = ADVANCE
var lineIndex = 0
var chunks = 30

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func getChunk(fileContent *[]string, from, to int) []string {
	return (*fileContent)[from:to]
}

/*
func NewVBox(c ...Widget) *Box {
	return &Box{
		children:  c,
		alignment: Vertical,
	}
}
*/

func clearBox(box *tui.Box, contentLen int) {
	for i := 0; i < contentLen; i++ {
		box.Append(tui.NewHBox(
			tui.NewLabel(""),
			tui.NewSpacer(),
		))
	}
}

func putText(box *tui.Box, content *[]string) {
	clearBox(box, len(*content))
	for _, txt := range *content {
		box.Append(tui.NewHBox(
			tui.NewLabel(txt),
			tui.NewSpacer(),
		))
	}
}

func main() {

	args := os.Args
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "error: missing file to read")
		os.Exit(1)
	}

	fileName := os.Args[1]
	fileContent, err := readLines(fileName)
	check(err)

	sidebar := tui.NewVBox(
		tui.NewLabel("CHANNELS"),
		tui.NewLabel("general"),
		tui.NewLabel("random"),
		tui.NewLabel(""),
		tui.NewLabel("DIRECT MESSAGES"),
		tui.NewLabel("slackbot"),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	history := tui.NewVBox()

	// All this code is totally optional ...
	// for _, m := range posts {
	// 	history.Append(tui.NewHBox(
	// 		// tui.NewLabel(m.line),
	// 		tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", m.line))),
	// 		//tui.NewLabel(m.message),
	// 		// tui.NewSpacer(),
	// 	))
	// }

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	// input.OnSubmit(func(e *tui.Entry) {
	// 	history.Append(tui.NewHBox(
	// 		tui.NewLabel(time.Now().Format("15:04")),
	// 		tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", "john"))),
	// 		tui.NewLabel(e.Text()),
	// 		tui.NewSpacer(),
	// 	))
	// 	input.SetText("")
	// })

	someChunk := getChunk(&fileContent, from, to)
	for _, txt := range someChunk {
		history.Append(tui.NewHBox(
			tui.NewLabel(txt),
			tui.NewSpacer(),
		))
	}

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	// down ...
	ui.SetKeybinding("j", func() {
		if from < len(fileContent) {
			from++
		}

		if from >= len(fileContent) {
			ui.Quit()
			fmt.Println("We are done ... ")
			return
		}

		if to < len(fileContent) {
			to++
		}
		//fmt.Printf("[%d, %d]", from, to)
		// fmt.Printf("(%s)\n", getChunk(&fileContent, from, to))
		chunk := getChunk(&fileContent, from, to)
		putText(history, &chunk)
	})

	// up ...
	ui.SetKeybinding("k", func() {
		fmt.Println("up ...")
	})

	ui.SetKeybinding("Esc", func() {
		ui.Quit()
		fmt.Println("Bye 2 ... ")
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
