package streamdeck

import (
	"fmt"
)

type Logger struct {
	conn *Conn
}

func (l *Logger) Println(a ...interface{}) error {
	return l.conn.Send(&LogMessage{
		Message: fmt.Sprintln(a...),
	})
}

func (l *Logger) Printf(format string, a ...interface{}) error {
	return l.conn.Send(&LogMessage{
		Message: fmt.Sprintf(format, a...),
	})
}
