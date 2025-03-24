package display

import (
	"fmt"
	"strings"

	"github.com/landru29/cnc-serial/internal/model"
)

// Write implements the io.Writer interface.
func (s *Screen) Write(data []byte) (int, error) {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	s.bufferData += string(data)

	splitter := strings.Split(s.bufferData, "\n")
	if len(splitter) < 2 { //nolint: mnd
		return len(data), nil
	}

	for _, line := range splitter {
		if status := model.DecodeStatus(line); status != nil {
			s.displayStatus(*status)

			continue
		}

		if response := model.DecodeResponse(line); response != nil {
			s.displayResponse(*response)

			continue
		}

		if request := model.DecodeRequest(line); request != nil {
			s.displayRequest(*request)

			continue
		}

		if program := model.DecodeProgram(line); program != nil {
			s.displayProgram(*program)

			continue
		}

		if line != "" {
			_, _ = s.logArea.Write([]byte(line + "\n"))
		}
	}

	s.bufferData = splitter[len(splitter)-1]

	return len(data), nil
}

func (s *Screen) displayStatus(status model.Status) {
	coordinates := status.ToolCoordinates()

	text := fmt.Sprintf(
		"%s\t\t%s\tX: %+07.2f\t\tY: %+07.2f\t\tZ: %+07.2f\t%04d\t%s",
		status.RelativeCoordinates,
		status.CurrentState(),
		coordinates.XCoordinate,
		coordinates.YCoordinate,
		coordinates.ZCoordinate,
		status.RemainingProgram,
		map[bool]string{true: "READY", false: "STOP "}[status.CanRun],
	)
	s.statusArea.SetText(text)

	userInputLabel := enterCommandLabel + " "

	if status.RemainingProgram != 0 {
		if status.CanRun {
			userInputLabel = enterCommandLabel + " ⌛ "
		} else {
			userInputLabel = enterCommandLabel + " ⚓ "
		}
	}

	s.userInput.SetLabel(userInputLabel)
}

func (s *Screen) displayResponse(status model.Response) {
	message := "[#00ff00]\t< " + status.Message + "[#505050]"

	if status.IsError {
		message = "<[#ff0000]\t< " + status.Message + "[#505050]"
	}

	_, _ = s.logArea.Write([]byte(message + "\n"))
}

func (s *Screen) displayRequest(status model.Request) {
	_, _ = s.logArea.Write([]byte("[#505050]> " + status.Message + "[#505050]" + "\n"))
}

func (s *Screen) displayProgram(program model.Program) {
	splitter := strings.Split(string(program.Data), "\n")
	output := make([]string, len(splitter))

	_, _, _, height := s.progArea.GetRect()

	scrollTo := 0

	if program.CurrentLine > height/2 {
		scrollTo = program.CurrentLine - height/3
	}

	for idx := range splitter {
		startColor := ""
		endColor := ""
		if idx == program.CurrentLine {
			startColor = "[#00ff00]"
			endColor = "[#c0c0c0]"
		}

		if idx < program.CurrentLine {
			startColor = "[#505050]"
		}

		output[idx] = fmt.Sprintf("[#888888]%d:[#ffffff] %s%s%s", idx, startColor, splitter[idx], endColor)
	}

	s.progArea.SetText(strings.Join(output, "\n"))

	s.progArea.ScrollTo(scrollTo, 0)
}
