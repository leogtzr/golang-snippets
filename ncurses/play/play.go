package main

import (
	"fmt"
	"log"
	"time"

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

var posts = []post{
	{line: "jane"},
}

func main() {
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

	for _, m := range posts {
		history.Append(tui.NewHBox(
			// tui.NewLabel(m.line),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", m.line))),
			//tui.NewLabel(m.message),
			tui.NewSpacer(),
		))
	}

	// fmt.Println(history.Alignment())

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

	input.OnSubmit(func(e *tui.Entry) {
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", "john"))),
			tui.NewLabel(e.Text()),
			tui.NewSpacer(),
		))
		input.SetText("")
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("j", func() {
		fmt.Println("Bye ... ")
		// ui.Quit()
		fmt.Println("Bye 2 ... ")
	})

	ui.SetKeybinding("Esc", func() {
		fmt.Println("Bye ... ")
		// ui.Quit()
		fmt.Println("Bye 2 ... ")
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
