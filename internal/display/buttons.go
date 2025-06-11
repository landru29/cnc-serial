//go:build withbutton

package display

import (
	"context"
	"strings"

	"github.com/landru29/cnc-serial/internal/control"
	"github.com/rivo/tview"
)

// Screen is the main layout.
type Screen struct {
	xButton axisButtons
	yButton axisButtons
	zButton axisButtons

	BaseScreen
}

type axisButtons struct {
	down *tview.Button
	up   *tview.Button
	name string
}

func (s *Screen) makeButtons(ctx context.Context) tview.Primitive { //nolint: ireturn
	s.xButton = s.newAxisButton(ctx, "x", "x→", "←x")
	s.yButton = s.newAxisButton(ctx, "y", "y↑", "↓y")
	s.zButton = s.newAxisButton(ctx, "z", "z↑", "↓z")

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

	return output
}

func (s *Screen) newAxisButton(ctx context.Context, name string, up string, down string) axisButtons {
	output := axisButtons{
		up:   tview.NewButton(up),
		down: tview.NewButton(down),
		name: strings.ToUpper(name),
	}

	output.up.SetBorder(true).SetRect(0, 0, 3, 3)   //nolint: mnd
	output.down.SetBorder(true).SetRect(0, 0, 3, 3) //nolint: mnd

	output.up.SetSelectedFunc(func() {
		output.move(ctx, s.commander, true, s.BaseScreen.navigationInc)
	})

	output.down.SetSelectedFunc(func() {
		output.move(ctx, s.commander, false, s.BaseScreen.navigationInc)
	})

	return output
}

func (a axisButtons) move(ctx context.Context, commander control.Commander, up bool, step float64) {
	if !up {
		step *= -1.0
	}

	_ = commander.MoveRelative(ctx, step, a.name)
}

func (s *Screen) layout(ctx context.Context) *tview.Flex {
	return tview.NewFlex().
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(s.statusArea, 1, 0, false).
				AddItem(tview.NewFlex().
					AddItem(s.logArea, 0, 1, false).
					AddItem(s.progArea, 0, 1, false),
									0, 8, false).
				AddItem(s.helpArea, 0, 2, false). //nolint: mnd
				AddItem(s.userInput, 0, 1, true), 0, 4, true).
		AddItem(
			s.makeButtons(ctx), 0, 1, false)
}
