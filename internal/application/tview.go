package application

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Display is the tview display.
func (c Client) Display() *tview.Application {
	return c.display
}

func (c *Client) buildView() {
	c.display = tview.NewApplication()

	c.logArea = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			c.display.Draw()
		})

	c.helpArea = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			c.display.Draw()
		})

	c.userInput = tview.NewInputField().
		SetLabel("Enter command ").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetDoneFunc(func(_ tcell.Key) {
			text := c.userInput.GetText()

			if text == "exit" {
				c.display.Stop()

				return
			}

			c.SendCommand(text)

			c.userInput.SetText("")
		})

	c.userInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyUp && c.cursor >= 0 {
			c.userInput.SetText(c.lastCommand[c.cursor])

			c.cursor--
		}

		if event.Key() == tcell.KeyDown && c.cursor >= 0 && c.cursor < len(c.lastCommand)-1 {
			c.userInput.SetText(c.lastCommand[c.cursor])

			c.cursor++
		}

		return event
	})

	c.userInput.SetChangedFunc(func(text string) {
		entry := strings.Split(text, " ")

		if description := c.codeDescription(strings.ToUpper(entry[0])); description != "" {
			c.helpArea.SetText(description)

			return
		}

		c.helpArea.SetText("")
	})

	c.logArea.SetBorder(true)
	c.helpArea.SetBorder(true)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(c.logArea, 0, 8, false).  //nolint: mnd
		AddItem(c.helpArea, 0, 2, false). //nolint: mnd
		AddItem(c.userInput, 0, 1, true)

	c.display = c.display.SetRoot(flex, true)
}

// Start launches the tview application.
func (c Client) Start() error {
	return c.display.Run()
}
