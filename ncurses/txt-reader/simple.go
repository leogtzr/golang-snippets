package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcusolsson/tui-go"
)

// ADVANCE ...
const Advance int = 30
const WrapMax = 80
const GotoWidgetIndex = 2

var from = 0
var to = Advance
var gotoLine = ""
var fileToOpen = flag.String("file", "", "File to open")
var openLatestFile = flag.Bool("latest", false, "Open the latest text file")

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

func downText(fileContent *[]string, txtArea *tui.Box) {
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
	putText(txtArea, &chunk)
}

func upText(fileContent *[]string, txtArea *tui.Box) {
	if from <= 0 {
		return
	}

	if from > 0 {
		from--
	}

	to--

	chunk := getChunk(fileContent, from, to)
	putText(txtArea, &chunk)
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
	return fmt.Sprintf("%s", strings.Join(fields, strings.Repeat(" ", wrapLength/(numberOfWords-1))))
}

func getStatusInformation(fileContent *[]string) string {
	return fmt.Sprintf("%d of %d%%                                   ", to, len(*fileContent))
}

func addUpBinding(fileContent *[]string, box *tui.Box, input *tui.Entry) func() {
	return func() {
		upText(fileContent, box)
		input.SetText(getStatusInformation(fileContent))
	}
}

func addDownBinding(fileContent *[]string, box *tui.Box, input *tui.Entry) func() {
	return func() {
		downText(fileContent, box)
		input.SetText(getStatusInformation(fileContent))
	}
}

func addGotoWidget(box *tui.Box) {
	gotoInput := tui.NewTextEdit()
	gotoInput.SetText("Go to line: ")
	gotoInput.SetFocused(true)
	gotoInput.OnTextChanged(func(entry *tui.TextEdit) {
		gotoLine = entry.Text()
	})
	box.Append(gotoInput)
}

func exists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

func saveStatus(fileName string, from, to int) {
	// write from, to y el nombre del archivo ...
	homeDir := os.Getenv("HOME")
	f, err := os.Create(filepath.Join(homeDir, "txtread"))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}
	defer f.Close()
	// w := bufio.NewWriter(f)
	// w.WriteString(fmt.Sprintf("%s|%d|%d", fileName, from, to))
	f.WriteString(fmt.Sprintf("%s|%d|%d", fileName, from, to))
}

func getFileNameFromLatest() (string, error) {
	homeDir := os.Getenv("HOME")
	latestFilePath := filepath.Join(homeDir, "txtread")

	if !exists(latestFilePath) {
		return "", fmt.Errorf("'%s' does not exist", latestFilePath)
	}

	f, err := os.Open(latestFilePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fileContent, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	latestFileFields := strings.Split(string(fileContent), "|")
	if len(latestFileFields) != 3 {
		return "", fmt.Errorf("Wrong format in '%s'", latestFilePath)
	}
	return latestFileFields[0], nil
}

func main() {

	flag.Parse()
	fileName := *fileToOpen
	if fileName == "" && !*openLatestFile {
		fmt.Fprintln(os.Stderr, "error: missing file to read")
		os.Exit(1)
	}

	if fileName != "" && *openLatestFile {
		fmt.Fprintln(os.Stderr, "error: conflicting options")
		os.Exit(1)
	}

	var err error

	if *openLatestFile {
		fileName, err = getFileNameFromLatest()
		if err != nil {
			log.Fatal(err)
		}
	}
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

	txtArea := tui.NewVBox()
	txtAreaScroll := tui.NewScrollArea(txtArea)
	txtAreaScroll.SetAutoscrollToBottom(true)

	txtAreaBox := tui.NewVBox(txtAreaScroll)
	txtAreaBox.SetBorder(true)

	inputCommand := tui.NewEntry()
	inputCommand.SetFocused(true)
	inputCommand.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputCommand.SetEchoMode(tui.EchoModeNormal)

	inputCommandBox := tui.NewHBox(inputCommand)
	inputCommandBox.SetBorder(true)
	inputCommandBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	txtReader := tui.NewVBox(txtAreaBox, inputCommandBox)
	txtReader.SetSizePolicy(tui.Expanding, tui.Expanding)

	chunk := getChunk(&fileContent, from, to)
	putText(txtArea, &chunk)

	root := tui.NewHBox(txtReader)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	// down ...
	ui.SetKeybinding("j", addDownBinding(&fileContent, txtArea, inputCommand))
	ui.SetKeybinding("Down", addDownBinding(&fileContent, txtArea, inputCommand))

	// Up ...
	ui.SetKeybinding("k", addUpBinding(&fileContent, txtArea, inputCommand))
	ui.SetKeybinding("Up", addUpBinding(&fileContent, txtArea, inputCommand))

	// go to:
	ui.SetKeybinding("g", func() {
		addGotoWidget(txtReader)
	})

	ui.SetKeybinding("r", func() {
		txtReader.Remove(GotoWidgetIndex)
		inputCommand.SetText(getStatusInformation(&fileContent))
	})

	ui.SetKeybinding("s", func() {
		// save status ...
		absoluteFilePath, _ := filepath.Abs(fileName)
		saveStatus(absoluteFilePath, from, to)
	})

	ui.SetKeybinding("Esc", func() {
		ui.Quit()
		fmt.Printf("Got in the buffer: [%s]\n", gotoLine)
	})

	inputCommand.SetText(getStatusInformation(&fileContent))

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
