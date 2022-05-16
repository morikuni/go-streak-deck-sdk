package main

import (
	"context"
	"fmt"
	"time"

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

	err = sdk.Receive(context.Background(), streamdeck.HandlerFunc(func(ctx context.Context, ev streamdeck.Event) error {
		switch ev := ev.(type) {
		case *streamdeck.KeyDown:
			sdk.Log("key down")
			return sdk.SetTitle(ev.Context, time.Now().Format("15:04:05"), streamdeck.TargetBoth, 0)
		case *streamdeck.KeyUp:
			sdk.Log("key up")
			return sdk.SetTitle(ev.Context, "", streamdeck.TargetBoth, 0)
		default:
			fmt.Println(ev)
			return nil
		}
	}))
	if err != nil {
		sdk.Log(err)
	}
}
