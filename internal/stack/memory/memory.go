// Package memory is the memory stack.
package memory

// Stack iss the command stack.
type Stack struct {
	content      []string
	cursor       int
	currentGiven bool
}

// New creates a new stack.
func New() *Stack {
	return &Stack{}
}

// Push add command in the stack.
func (s *Stack) Push(cmd string) {
	s.content = append([]string{cmd}, s.content...)
}

// ResetCursor set the cursor on the last command.
func (s *Stack) ResetCursor() {
	s.cursor = 0
	s.currentGiven = false
}

// NavigateUp navigates to previous command.
func (s *Stack) NavigateUp() string {
	if len(s.content) == 0 {
		return ""
	}

	if !s.currentGiven {
		s.currentGiven = true

		return s.content[0]
	}

	if s.cursor < len(s.content)-1 {
		s.cursor++

		return s.content[s.cursor]
	}

	return ""
}

// NavigateDown navigates to next command.
func (s *Stack) NavigateDown() string {
	if len(s.content) == 0 {
		return ""
	}

	if !s.currentGiven {
		s.currentGiven = true

		return s.content[0]
	}

	if s.cursor > 0 {
		s.cursor--

		return s.content[s.cursor]
	}

	return ""
}
