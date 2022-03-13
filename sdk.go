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

func (sdk *SDK) SetTitle(context InstanceID, title string, target Target, state int) error {
	return sdk.conn.Send(&SetTitle{
		Context: context,
		Title:   title,
		Target:  target,
		State:   state,
	})
}

func (sdk *SDK) SetImage(context InstanceID, img Image, target Target, state int) error {
	return sdk.conn.Send(&SetImage{
		Context: context,
		Image:   img,
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

func (sdk *SDK) Receive(ctx context.Context, h Handler) error {
	hc := &handlerCtx{
		ctx,
		sdk,
	}

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

		err = h.Handle(hc, ev)
		if err != nil {
			return err
		}
	}
}

type Context interface {
	context.Context

	SDK() *SDK
	OpenURL(url string) error
	Log(a ...interface{})
	Logf(format string, a ...interface{})
}

type handlerCtx struct {
	context.Context

	sdk *SDK
}

var _ Context = (*handlerCtx)(nil)

func (ctx *handlerCtx) SDK() *SDK {
	return ctx.sdk
}

func (ctx *handlerCtx) OpenURL(url string) error {
	return ctx.sdk.OpenURL(url)
}

func (ctx *handlerCtx) Log(a ...interface{}) {
	ctx.sdk.Log(a...)
}

func (ctx *handlerCtx) Logf(format string, a ...interface{}) {
	ctx.sdk.Logf(format, a...)
}

type Handler interface {
	Handle(ctx Context, ev Event) error
}

type HandlerFunc func(ctx Context, ev Event) error

func (f HandlerFunc) Handle(ctx Context, ev Event) error {
	return f(ctx, ev)
}
