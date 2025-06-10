package gpm

import "github.com/gdamore/tcell/v2"

func NewScreen() (tcell.Screen, error) {
	// Windows is happier if we try for a console screen first.
	if s, _ := tcell.NewConsoleScreen(); s != nil {
		return s, nil
	} else if s, e := tcell.NewTerminfoScreen(); s != nil {
		s, _ = WarpGPMSupport(s)
		return s, nil
	} else {
		return nil, e
	}
}
