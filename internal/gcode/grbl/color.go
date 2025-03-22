package grbl

import (
	"strings"
)

// Colorize applies the syntax highlighting.
func (g Gerbil) Colorize(text string) string {
	splitter := strings.Split(strings.ToUpper(strings.TrimSpace(text)), " ")
	newSplitter := make([]string, len(splitter))

	for idx, word := range splitter {
		switch {
		case strings.HasPrefix(word, "G"):
			newSplitter[idx] = "[#fffc3a]" + word
		case strings.HasPrefix(word, "M"):
			newSplitter[idx] = "[#ffaff8]" + word
		case strings.HasPrefix(word, "X"):
			newSplitter[idx] = "[#00ff00]" + word
		case strings.HasPrefix(word, "Y"):
			newSplitter[idx] = "[#69f1ff]" + word
		case strings.HasPrefix(word, "Z"):
			newSplitter[idx] = "[#ff0000]" + word
		default:
			newSplitter[idx] = "[white]" + word
		}
	}

	return strings.Join(newSplitter, " ")
}
