package display

import (
	"fmt"
	"strings"

	"github.com/landru29/serial/internal/control"
	"github.com/rivo/tview"
)

type axisButtons struct {
	down *tview.Button
	up   *tview.Button
	name string
}

func (s *Screen) makeButtons() tview.Primitive { //nolint: ireturn
	s.xButton = s.newAxisButton("x", "←x", "x→")
	s.yButton = s.newAxisButton("y", "y↑", "↓y")
	s.zButton = s.newAxisButton("z", "z↑", "↓z")

	xyGrid := tview.NewGrid().
		SetRows(0, 0, 0).
		SetColumns(0, 0, 0).
		AddItem(s.yButton.up, 0, 1, 1, 1, 0, 0, false).
		AddItem(s.yButton.down, 2, 1, 1, 1, 0, 0, false). //nolint: mnd
		AddItem(s.xButton.down, 1, 0, 1, 1, 0, 0, false).
		AddItem(s.xButton.up, 1, 2, 1, 1, 0, 0, false) //nolint: mnd

	xyGrid.SetBorder(true)

	zGrid := tview.NewGrid().
		SetRows(0, 0).
		SetColumns(0, 0, 0).
		AddItem(s.zButton.up, 0, 1, 1, 1, 0, 0, false).
		AddItem(s.zButton.down, 1, 1, 1, 1, 0, 0, false)

	zGrid.SetBorder(true)

	output := tview.NewGrid().SetRows(0, 0).SetColumns(0).
		AddItem(xyGrid, 0, 0, 1, 1, 0, 0, false).
		AddItem(zGrid, 1, 0, 1, 1, 0, 0, false)

	output.SetBorder(true)

	return output
}

func (s *Screen) newAxisButton(name string, up string, down string) axisButtons {
	output := axisButtons{
		up:   tview.NewButton(up),
		down: tview.NewButton(down),
		name: strings.ToUpper(name),
	}

	output.up.SetBorder(true).SetRect(0, 0, 3, 3)   //nolint: mnd
	output.down.SetBorder(true).SetRect(0, 0, 3, 3) //nolint: mnd

	output.up.SetSelectedFunc(func() {
		output.move(s.commander, true)
	})

	output.down.SetSelectedFunc(func() {
		output.move(s.commander, true)
	})

	return output
}

func (a axisButtons) move(commander control.Commander, up bool) {
	commands := []string{
		"G91",
		fmt.Sprintf("G1 %s%s10", a.name, map[bool]string{true: "+", false: "-"}[up]),
		"G90",
	}

	if commander.IsRelative() {
		commands = commands[1:2]
	}

	_ = commander.Send(commands...)
}
