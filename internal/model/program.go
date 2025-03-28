package model

import (
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Program is a program with execution cursor.
type Program struct {
	Data        []byte
	CurrentLine int64
}

// Encode is the program encoder.
func (p *Program) Encode(writer io.Writer) error {
	if p == nil {
		return nil
	}

	if _, err := fmt.Fprintf(writer, "%d|%s\n", p.CurrentLine, base64.StdEncoding.EncodeToString(p.Data)); err != nil {
		return err
	}

	return nil
}

// DecodeProgram is the program decoder.
func DecodeProgram(data string) *Program {
	splitter := strings.Split(data, "|")
	if len(splitter) != 2 { //nolint: mnd
		return nil
	}

	currentLine, err := strconv.ParseInt(splitter[0], 10, 64)
	if err != nil {
		return nil
	}

	prog, err := base64.StdEncoding.DecodeString(splitter[1])
	if err != nil {
		return nil
	}

	return &Program{
		Data:        prog,
		CurrentLine: currentLine,
	}
}
