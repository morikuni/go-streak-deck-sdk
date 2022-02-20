package streamdeck

import (
	"context"
	"errors"
	"fmt"
	"io"
	"runtime/debug"
)

type SDK struct {
	conn *Conn
}

func NewSDK(conn *Conn) *SDK {
	return &SDK{conn: conn}
}

func (sdk *SDK) ShowOK(context InstanceID) error {
	return sdk.conn.Send(&ShowOK{
		Context: context,
	})
}

func (sdk *SDK) ShowAlert(context InstanceID) error {
	return sdk.conn.Send(&ShowAlert{
		Context: context,
	})
}

// Log prints log via Stream Deck API.
func (sdk *SDK) Log(a ...interface{}) {
	s := fmt.Sprintln(a...)
	_ = sdk.conn.Send(&LogMessage{
		Message: s[:len(s)-1],
	})
}

// Logf prints formatted log via Stream Deck API.
func (sdk *SDK) Logf(format string, a ...interface{}) {
	_ = sdk.conn.Send(&LogMessage{
		Message: fmt.Sprintf(format, a...),
	})
}

func (sdk *SDK) WatchInstance(ctx context.Context, f InstanceFactory) error {
	defer func() {
		// TODO: remove recover
		if r := recover(); r != nil {
			sdk.Log(r)
			sdk.Log(string(debug.Stack()))
		}
	}()

	s := newSupervisor(ctx, sdk, f)

	for {
		ev, err := sdk.conn.Receive()
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("stop due to EOF: %w", err)
		}
		if err != nil {
			sdk.Log("go-stream-deck-sdk: error on receive", err)
			continue
		}

		s.handle(ev)
	}
}
