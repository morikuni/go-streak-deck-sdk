package main

import (
	"time"

	"golang.org/x/net/context"

	streamdeck "github.com/morikuni/go-stream-deck-sdk"
)

func main() {
	conn, err := streamdeck.Dial()
	if err != nil {
		// This log will be discarded because logs should be sent to the stream deck as a command.
		panic(err)
	}
	defer conn.Close()

	sdk := streamdeck.NewSDK(conn)
	sdk.Log("start")
	defer func() {
		sdk.Log("exit", recover())
	}()

	factory := func(ctx streamdeck.InstanceContext, id streamdeck.InstanceID) *streamdeck.Instance {
		return &streamdeck.Instance{
			OnKeyDown: func(ctx streamdeck.InstanceContext, ev *streamdeck.KeyDown) error {
				ctx.Log("key down")
				return ctx.SetTitle(time.Now().Format("15:04:05"), streamdeck.TitleTargetBoth, 0)
			},
			OnKeyUp: func(ctx streamdeck.InstanceContext, ev *streamdeck.KeyUp) error {
				ctx.Log("key up")
				return ctx.SetTitle("", streamdeck.TitleTargetBoth, 0)
			},
		}
	}
	mux := streamdeck.NewMux(factory)

	err = sdk.Receive(context.Background(), mux)
	if err != nil {
		sdk.Log(err)
	}
}
