package main

import (
	"log"

	"github.com/marcusolsson/tui-go"
)

func main() {
	buffer := tui.NewLabel(body)
	buffer.SetSizePolicy(tui.Expanding, tui.Expanding)
	buffer.SetText(body)
	buffer.SetFocused(true)
	// buffer.SetWordWrap(true)

	status := tui.NewStatusBar("lorem.txt")
	status.SetText("cmn")

	root := tui.NewVBox(buffer, status)
	root.SetBorder(true)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

const body = `
Maecenas eget tristique dolor. Quisque vel velit ante. Pellentesque habitant morbi tristique senectus et netus et 
\tcmn
`
