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
		"%s\t\t%s\tX: %+07.2f\t\tY: %+07.2f\t\tZ: %+07.2f",
		status.RelativeCoordinates,
		status.CurrentState(),
		coordinates.XCoordinate,
		coordinates.YCoordinate,
		coordinates.ZCoordinate,
	)
	s.statusArea.SetText(text)
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

	_, _, _, height := s.progArea.GetRect()

	scrollTo := 0

	if program.CurrentLine > height/2 {
		scrollTo = program.CurrentLine - height/2
	}

	if program.CurrentLine <= len(splitter) {
		out := splitter[:program.CurrentLine]
		out = append(out, "[#00ff00]"+splitter[program.CurrentLine]+"[#c0c0c0]")
		splitter = append(out, splitter[program.CurrentLine+1:]...)
	}

	s.progArea.SetText("[#505050]" + strings.Join(splitter, "\n"))

	s.progArea.ScrollTo(scrollTo, 0)
}
