// Package model gathers all data models.
package model

import (
	"bytes"
	"encoding/json"
	"io"
)

func encode(writer io.Writer, object any) error {
	var buffer bytes.Buffer

	if err := json.NewEncoder(&buffer).Encode(object); err != nil {
		return err
	}

	if err := buffer.WriteByte('\n'); err != nil {
		return err
	}

	if _, err := io.Copy(writer, &buffer); err != nil {
		return err
	}

	return nil
}

// ObjectName is an object name.
type ObjectName string
