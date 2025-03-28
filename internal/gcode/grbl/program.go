package grbl

import (
	"bytes"
	"errors"
	"io"
	"strings"

	apperrors "github.com/landru29/cnc-serial/internal/errors"
	"github.com/landru29/cnc-serial/internal/gcode"
	"github.com/landru29/cnc-serial/internal/model"
)

var _ gcode.Programmer = &Program{}

// Program is a program.
type Program struct {
	content             []byte
	buffer              *bytes.Buffer
	carriageReturnCount int64
	currentCommand      string
	commandsToBeRead    int64
}

// NewProgram creates the program.
func NewProgram(data io.Reader) (*Program, error) {
	content, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}

	return &Program{
		content: content,
		buffer:  bytes.NewBuffer(content),
	}, nil
}

// Reset implements the gcode.Programmer interface.
func (p *Program) Reset() error {
	p.buffer = bytes.NewBuffer(p.content)

	p.carriageReturnCount = 0

	p.currentCommand = ""

	p.commandsToBeRead = 0

	return p.ReadNextInstruction()
}

// SetLinesToExecute implements the gcode.Programmer interface.
func (p *Program) SetLinesToExecute(count int64) {
	p.commandsToBeRead = count
}

// NextCommandToExecute implements the gcode.Programmer interface.
func (p *Program) NextCommandToExecute() (string, error) {
	defer func() {
		if p.commandsToBeRead > 0 {
			p.commandsToBeRead--
		}
	}()

	if p.commandsToBeRead != 0 {
		cmdToSend := p.currentCommand

		return cmdToSend, p.ReadNextInstruction()
	}

	return "", apperrors.ErrProgramIdle
}

// ToModel implements the gcode.Programmer interface.
func (p *Program) ToModel() *model.Program {
	if p == nil {
		return nil
	}

	return &model.Program{
		Data:        p.content,
		CurrentLine: p.CurrentLine(),
	}
}

// CurrentLine implements the gcode.Programmer interface.
func (p Program) CurrentLine() int64 {
	return p.carriageReturnCount
}

// CurrentCommand implements the gcode.Programmer interface.
func (p Program) CurrentCommand() string {
	return p.currentCommand
}

// Content implements the gcode.Programmer interface..
func (p Program) Content() string {
	return string(p.content)
}

// ReadNextInstruction implements the gcode.Programmer interface.
func (p *Program) ReadNextInstruction() error {
	if p == nil || p.buffer == nil {
		return nil
	}

	if err := p.skipComments(); err != nil {
		return err
	}

	line := p.readNextChars(false, '\n', ';', '(')

	p.currentCommand = strings.TrimSpace(line)

	return nil
}

func (p *Program) skipSpaces() {
	for {
		current, err := p.buffer.ReadByte()
		if errors.Is(err, io.EOF) {
			return
		}

		if current == '\n' {
			p.carriageReturnCount++
		}

		if current != ' ' && current != '\t' && current != '\r' && current != '\n' {
			break
		}
	}

	_ = p.buffer.UnreadByte()
}

func (p *Program) readNextChars(skipStopper bool, stoppers ...byte) string {
	output := ""

	var lastRead byte

	for {
		current, err := p.buffer.ReadByte()
		if errors.Is(err, io.EOF) {
			return ""
		}

		output += string(current)

		quit := false

		for _, character := range stoppers {
			if current == character {
				quit = true

				break
			}
		}

		if current == '\n' {
			p.carriageReturnCount++
		}

		if quit {
			lastRead = current

			break
		}
	}

	if !skipStopper {
		if lastRead == '\n' {
			p.carriageReturnCount--
		}

		_ = p.buffer.UnreadByte()
	}

	return output[:len(output)-1]
}

func (p *Program) skipComments() error {
	// skip spaces.
	p.skipSpaces()

	// detect comment marker.
	current, err := p.buffer.ReadByte()
	if errors.Is(err, io.EOF) {
		return nil
	}

	// Is there a comment ?
	if current != '(' && current != '%' && current != ';' {
		err := p.buffer.UnreadByte()
		if errors.Is(err, io.EOF) || err == nil {
			return nil
		}

		return err
	}

	// comment on one line.
	if current == '%' || current == ';' {
		_ = p.readNextChars(true, '\n')

		return p.skipComments()
	}

	// multiline comment.
	_ = p.readNextChars(true, ')')

	return p.skipComments()
}
