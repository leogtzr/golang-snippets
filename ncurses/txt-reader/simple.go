package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/marcusolsson/tui-go"
)

// ADVANCE ...
const Advance int = 30
const WrapMax = 80

var from = 0
var to = Advance

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
		txt = strings.Replace(txt, " ", " ", -1)
		txt = strings.Replace(txt, "\t", "    ", -1)
		txt = wrap(txt)
		box.Append(tui.NewVBox(
			tui.NewLabel(txt),
			tui.NewSpacer(),
		))
	}
}

func downText(fileContent *[]string, history *tui.Box) {
	if from < len(*fileContent) {
		from++
	}
	if to >= len(*fileContent) {
		return
	}

	if to < len(*fileContent) {
		to++
	}
	chunk := getChunk(fileContent, from, to)
	putText(history, &chunk)
}

func upText(fileContent *[]string, history *tui.Box) {
	if from <= 0 {
		return
	}

	if from > 0 {
		from--
	}

	to--

	chunk := getChunk(fileContent, from, to)
	putText(history, &chunk)
}

func needsSemiWrap(line string) bool {
	len := len(line)
	if len < (WrapMax / 2) {
		return false
	}
	return (len > (WrapMax / 2)) && (len < WrapMax)
}

func countWithoutWhitespaces(words []string) int {
	count := 0
	for _, w := range words {
		count += len(w)
	}
	return count
}

func wrap(line string) string {
	if !needsSemiWrap(line) {
		return line
	}
	fields := strings.Fields(line)
	numberOfWords := len(fields)
	countWithoutSpaces := countWithoutWhitespaces(fields)
	wrapLength := WrapMax - countWithoutSpaces
	if numberOfWords == 1 || numberOfWords == 0 {
		return line
	}
	return fmt.Sprintf("<%s>", strings.Join(fields, strings.Repeat(" ", wrapLength/(numberOfWords-1))))
}

func addUpBinding(fileContent *[]string, box *tui.Box, input *tui.Entry) func() {
	return func() {
		upText(fileContent, box)
		input.SetText(fmt.Sprintf("%d of %d%%                                   ", to, len(*fileContent)))
	}
}

func addDownBinding(fileContent *[]string, box *tui.Box, input *tui.Entry) func() {
	return func() {
		downText(fileContent, box)
		input.SetText(fmt.Sprintf("%d of %d%%                                   ", to, len(*fileContent)))
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

	// sidebar := tui.NewVBox(
	// 	tui.NewLabel("CHANNELS"),
	// 	tui.NewLabel("general"),
	// 	tui.NewLabel("random"),
	// 	tui.NewLabel(""),
	// 	tui.NewLabel("DIRECT MESSAGES"),
	// 	tui.NewLabel("slackbot"),
	// 	tui.NewSpacer(),
	// )
	// sidebar.SetBorder(true)

	history := tui.NewVBox()
	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)
	input.SetEchoMode(tui.EchoModeNormal)

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
	putText(history, &someChunk)

	//root := tui.NewHBox(sidebar, chat)
	root := tui.NewHBox(chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	// down ...
	ui.SetKeybinding("j", addDownBinding(&fileContent, history, input))
	ui.SetKeybinding("Down", addDownBinding(&fileContent, history, input))
	ui.SetKeybinding("Enter", addDownBinding(&fileContent, history, input))

	// Up ...
	ui.SetKeybinding("k", addUpBinding(&fileContent, history, input))
	ui.SetKeybinding("Up", addUpBinding(&fileContent, history, input))

	// go to:
	ui.SetKeybinding("g", func() {
		// root.Append
		gotoInput := tui.NewTextEdit()
		gotoInput.SetSizePolicy(tui.Expanding, tui.Expanding)
		gotoInput.SetText("Goto: ")
		gotoInput.SetFocused(true)
		gotoInput.SetWordWrap(true)
		chat.Append(gotoInput)
	})

	ui.SetKeybinding("r", func() {
		// fmt.Println(chat.Length())
		chat.Remove(2)
	})

	ui.SetKeybinding("Esc", func() {
		ui.Quit()
		fmt.Println("Bye 2 ... ")
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
