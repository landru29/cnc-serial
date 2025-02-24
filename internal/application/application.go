// Package application is the main application.
package application

import (
	"fmt"
	"strings"

	"github.com/landru29/serial/internal/gcode"
	"github.com/rivo/tview"
	"go.bug.st/serial"
)

// Client is the main application structure.
type Client struct {
	lastCommand  []string
	cursor       int
	port         serial.Port
	display      *tview.Application
	userInput    *tview.InputField
	logArea      *tview.TextView
	helpArea     *tview.TextView
	translations map[gcode.Language]gcode.CodeSet
	language     gcode.Language
}

// NewClient initializes a new application client.
func NewClient() (*Client, error) {
	output := &Client{
		language: gcode.DefaultLanguage,
	}

	output.buildView()

	translations, err := gcode.ReadCodes()
	if err != nil {
		return nil, err
	}

	output.translations = translations

	return output, nil
}

// SetLanguage sets the language.
func (c *Client) SetLanguage(language gcode.Language) {
	c.language = language
}

func (c Client) dryRun() bool {
	return c.port == nil
}

func (c Client) codeDescription(code string) string {
	codeName := gcode.CodeName(strings.ToUpper(code))
	if description := c.translations[c.language][codeName].Description; description != "" {
		return fmt.Sprintf("%s - %s", codeName, description)
	}

	return ""
}
