package streamdeck

import (
	"fmt"
)

type Logger struct {
	conn *Conn
}

// Println prints log via Stream Deck API.
func (l *Logger) Println(a ...interface{}) error {
	s := fmt.Sprintln(a...)
	return l.conn.Send(&LogMessage{
		Message: s[:len(s)-1],
	})
}

// Printf prints formatted log via Stream Deck API.
func (l *Logger) Printf(format string, a ...interface{}) error {
	return l.conn.Send(&LogMessage{
		Message: fmt.Sprintf(format, a...),
	})
}
