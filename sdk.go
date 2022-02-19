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

func (sdk *SDK) ShowOK(context InstanceID) {
	_ = sdk.conn.Send(&ShowOK{
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
		if r := recover(); r != nil {
			sdk.Log(r)
			sdk.Log(string(debug.Stack()))
		}
	}()

	instanceByID := make(map[InstanceID]*Instance)
	withInstance := func(id InstanceID, callback func(*Instance) error) {
		instance, ok := instanceByID[id]
		if !ok {
			instance = f(&instanceCtx{ctx, sdk.conn, id}, id)
			if instance == nil {
				sdk.Log("go-stream-deck-sdk: no instance returned: id = %s", id)
				return
			}
			instance.id = id
			instanceByID[id] = instance
		}

		err := callback(instance)
		if err != nil {
			sdk.Logf("go-stream-deck-sdk: error on instance(id=%v): %v", instance.id, err)
			return
		}
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

		switch t := ev.(type) {
		case *KeyDown:
			withInstance(t.Context, func(instance *Instance) error {
				return instance.OnKeyDown(instance.ctx(ctx, sdk.conn), t)
			})
		}
	}
}
