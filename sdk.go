package streamdeck

import (
	"context"
	"errors"
	"fmt"
	"io"
)

type SDK struct {
	conn *Conn

	debugLog bool
}

func NewSDK(conn *Conn) *SDK {
	return &SDK{conn: conn, debugLog: true}
}

func (sdk *SDK) OpenURL(url string) error {
	return sdk.conn.Send(&OpenURL{
		URL: url,
	})
}

func (sdk *SDK) SetTitle(context InstanceID, title string, target TitleTarget, state int) error {
	return sdk.conn.Send(&SetTitle{
		Context: context,
		Title:   title,
		Target:  target,
		State:   state,
	})
}

func (sdk *SDK) ShowAlert(context InstanceID) error {
	return sdk.conn.Send(&ShowAlert{
		Context: context,
	})
}

func (sdk *SDK) ShowOK(context InstanceID) error {
	return sdk.conn.Send(&ShowOK{
		Context: context,
	})
}

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

func (sdk *SDK) debug(a ...interface{}) {
	if sdk.debugLog {
		sdk.Log(a...)
	}
}

func (sdk *SDK) debugf(format string, a ...interface{}) {
	if sdk.debugLog {
		sdk.Logf(format, a...)
	}
}

func (sdk *SDK) WatchInstance(ctx context.Context, f InstanceFactory) error {
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

		sdk.debugf("[DEBUG] go-stream-deck-sdk: received: %#v", ev)

		s.handle(ev)
	}
}

type Context interface {
	context.Context

	OpenURL(url string) error
	Log(a ...interface{})
	Logf(format string, a ...interface{})
}

type Handler interface {
	Handle(ctx Context, ev Event)
}

type HandlerFunc func(ctx Context, ev Event)

func (f HandlerFunc) Handle(ctx Context, ev Event) {
	f(ctx, ev)
}
