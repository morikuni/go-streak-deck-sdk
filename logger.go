package streamdeck

import (
	"fmt"
)

type Logger struct {
	conn *Conn
}

// Println prints log via Stream Deck API.
func (l *Logger) Println(a ...interface{}) error {
	return l.conn.Send(&LogMessage{
		Message: fmt.Sprint(a...),
	})
}

// Printf prints formatted log via Stream Deck API.
func (l *Logger) Printf(format string, a ...interface{}) error {
	return l.conn.Send(&LogMessage{
		Message: fmt.Sprintf(format, a...),
	})
}
