// Package stack manages the stack command
package stack

// Pusher is the stack writer.
type Pusher interface {
	Push(cmd string)
}

// Retriever is the stack reader.
type Retriever interface {
	ResetCursor()
	NavigateUp() string
	NavigateDown() string
}

// Stacker is the stack manager.
type Stacker interface {
	Pusher
	Retriever
}
