// Package lang is the language management.
package lang

// Language is the user language.
type Language string

// DefaultLanguage is the default language.
const DefaultLanguage Language = "en"

// String implements the pflag.Value interface.
func (l Language) String() string {
	return string(l)
}

// Set implements the pflag.Value interface.
func (l *Language) Set(data string) error {
	*l = Language(data)

	return nil
}

// Type implements the pflag.Value interface.
func (l Language) Type() string {
	return "lang"
}
