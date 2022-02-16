package streamdeck

import (
	"fmt"
)

type SDK struct {
	conn *Conn
}

func NewFramework(conn *Conn) *SDK {
	return &SDK{conn: conn}
}

func (sdk *SDK) ShowOK(context InstanceID) error {
	return sdk.conn.Send(&ShowOK{
		Context: context,
	})
}

// Log prints log via Stream Deck API.
func (sdk *SDK) Log(a ...interface{}) error {
	s := fmt.Sprintln(a...)
	return sdk.conn.Send(&LogMessage{
		Message: s[:len(s)-1],
	})
}

// Logf prints formatted log via Stream Deck API.
func (sdk *SDK) Logf(format string, a ...interface{}) error {
	return sdk.conn.Send(&LogMessage{
		Message: fmt.Sprintf(format, a...),
	})
}
