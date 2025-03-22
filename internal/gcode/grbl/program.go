package grbl

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/landru29/cnc-serial/internal/gcode"
	"github.com/landru29/cnc-serial/internal/model"
)

var _ gcode.Programmer = &Program{}

// Program is a program.
type Program struct {
	content             []byte
	buffer              *bytes.Buffer
	carriageReturnCount int
	currentCommand      string
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

// ToModel is the data converter.
func (p *Program) ToModel() *model.Program {
	if p == nil {
		return nil
	}

	return &model.Program{
		Data:        p.content,
		CurrentLine: p.CurrentLine(),
	}
}

// CurrentLine is the line of the current instruction.
func (p Program) CurrentLine() int {
	return p.carriageReturnCount
}

// CurrentCommand is the current selected command.
func (p Program) CurrentCommand() string {
	return p.currentCommand
}

// Content is the line of the current instruction.
func (p Program) Content() string {
	return string(p.content)
}

// ReadNextInstruction move the pointer to the next instruction and read it.
func (p *Program) ReadNextInstruction() (string, error) {
	if err := p.skipComments(); err != nil {
		return "", err
	}

	line := p.readNextChars(false, '\n', ';', '(')

	p.currentCommand = strings.TrimSpace(line)

	return p.currentCommand, nil
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

		p.skipComments()

		return nil
	}

	// multiline comment.
	_ = p.readNextChars(true, ')')

	return p.skipComments()
}
