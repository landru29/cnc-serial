// Package application is the main application.
package application

import (
	"fmt"
	"strings"

	"github.com/landru29/serial/internal/control"
	"github.com/landru29/serial/internal/control/nop"
	"github.com/landru29/serial/internal/display"
	"github.com/landru29/serial/internal/gcode"
)

// Client is the main application structure.
type Client struct {
	commander    control.Commander
	screen       *display.Screen
	translations map[gcode.Language]gcode.CodeSet
	language     gcode.Language
}

// NewClient initializes a new application client.
func NewClient() (*Client, error) {
	output := &Client{
		language: gcode.DefaultLanguage,
	}

	output.screen = display.New(func(str string) string {
		return output.codeDescription(str)
	})

	output.commander = nop.New(output.screen.Output())

	output.screen.SetCommander(output.commander)

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

func (c Client) codeDescription(code string) string {
	codeName := gcode.CodeName(strings.ToUpper(code))
	if description := c.translations[c.language][codeName].Description; description != "" {
		return fmt.Sprintf("%s - %s", codeName, description)
	}

	return ""
}

// AvailableLanguages lists all available languages.
func (c Client) AvailableLanguages() []string {
	output := []string{}

	for lang := range c.translations {
		output = append(output, string(lang))
	}

	return output
}

// Start launches the tview application.
func (c Client) Start() error {
	return c.screen.Start()
}

// Bind is the main application loop.
func (c Client) Bind() {
	c.commander.Bind()
}
