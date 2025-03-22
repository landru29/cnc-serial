// Package grbl manages GERBIL specificities.
package grbl

import (
	"embed"
	"encoding/json"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/landru29/cnc-serial/internal/gcode"
	"github.com/landru29/cnc-serial/internal/lang"
)

//go:embed lang/*.json
var dataFS embed.FS

var _ gcode.Processor = &Gerbil{}

// Gerbil is the GERBIL standard.
type Gerbil struct {
	helper         map[lang.Language]CodeSet
	stateRegexp    *regexp.Regexp
	argumentRegexp *regexp.Regexp
	numberRegexp   *regexp.Regexp
}

// CodeDescription implements the Helper interface.
func (g Gerbil) CodeDescription(lang lang.Language, code string) string {
	if g.helper[lang] == nil {
		return ""
	}

	codeName := strings.ToUpper(strings.TrimSpace(code))

	if code, found := g.helper[lang][CodeName(codeName)]; found {
		return codeName + " - " + code.Description
	}

	return ""
}

// CodeName is a code name.
type CodeName string

// Code is a gcode instruction.
type Code struct {
	Description string `json:"description"`
}

// CodeSet is a set of codes.
type CodeSet map[CodeName]Code

// New reads and parse json data.
func New() (*Gerbil, error) {
	output := Gerbil{
		helper: map[lang.Language]CodeSet{},
	}

	stateRegexp, err := regexp.Compile(`^([a-zA-Z]+)`) //nolint: gocritic
	if err != nil {
		return nil, err
	}

	output.stateRegexp = stateRegexp

	argumentRegexp, err := regexp.Compile(`^([a-zA-Z]+)`) //nolint: gocritic
	if err != nil {
		return nil, err
	}

	output.argumentRegexp = argumentRegexp

	numberRegexp, err := regexp.Compile(`^([+-]?[0-9]+(\.[0-9]+)?)`) //nolint: gocritic
	if err != nil {
		return nil, err
	}

	output.numberRegexp = numberRegexp

	if errWalk := fs.WalkDir(dataFS, ".", func(pathName string, entry fs.DirEntry, pathErr error) error {
		if pathErr != nil {
			return pathErr
		}

		if entry.IsDir() {
			return nil
		}

		if strings.ToLower(filepath.Ext(pathName)) != ".json" {
			return nil
		}

		file, err := dataFS.Open(pathName)
		if err != nil {
			return err
		}

		defer func() {
			_ = file.Close()
		}()

		translation := CodeSet{}

		if err := json.NewDecoder(file).Decode(&translation); err != nil {
			return err
		}

		countryName := strings.Split(filepath.Base(pathName), ".")

		output.helper[lang.Language(strings.ToLower(countryName[0]))] = translation

		return nil
	}); errWalk != nil {
		return nil, errWalk
	}

	return &output, nil
}

// AvailableLanguages implements the gcode.Processor interface.
func (g Gerbil) AvailableLanguages() []lang.Language {
	output := []lang.Language{}
	for key := range g.helper {
		output = append(output, key)
	}

	return output
}
