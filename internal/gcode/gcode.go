// Package gcode manages the g-code.
package gcode

import (
	"embed"
	"encoding/json"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed *.json
var dataFS embed.FS

// CodeName is a code name.
type CodeName string

// Code is a gcode instruction.
type Code struct {
	Description string `json:"description"`
}

// CodeSet is a set of codes.
type CodeSet map[CodeName]Code

// ReadCodes reads and parse json data.
func ReadCodes() (map[Language]CodeSet, error) {
	output := map[Language]CodeSet{}

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

		output[Language(strings.ToLower(countryName[0]))] = translation

		return nil
	}); errWalk != nil {
		return nil, errWalk
	}

	return output, nil
}
