package model

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Program struct {
	Data        []byte
	CurrentLine int
}

// Encode is the program encoder.
func (p *Program) Encode(writer io.Writer) error {
	if p == nil {
		return nil
	}

	fmt.Fprintf(writer, "%d|%s\n", p.CurrentLine, base64.StdEncoding.EncodeToString(p.Data))
	return nil
}

// DecodeRequest is the program decoder.
func DecodeProgram(data string) *Program {
	splitter := strings.Split(data, "|")
	if len(splitter) != 2 {
		return nil
	}

	currentLine, err := strconv.ParseInt(splitter[0], 10, 64)
	if err != nil {
		return nil
	}

	os.WriteFile("/tmp/foo.txt", []byte(splitter[1]), os.ModePerm)

	prog, err := base64.StdEncoding.DecodeString(splitter[1])
	if err != nil {
		return nil
	}

	return &Program{
		Data:        prog,
		CurrentLine: int(currentLine),
	}
}
