package streamdeck

import (
	"golang.org/x/net/context"
)

type Context interface {
	context.Context

	ShowOK() error
}

type instanceCtx struct {
	context.Context

	conn       *Conn
	instanceID InstanceID
}

func (ctx *instanceCtx) ShowOK() error {
	return ctx.conn.Send(&ShowOK{
		Context: ctx.instanceID,
	})
}

type InstanceFactory func(ctx Context, id InstanceID) *Instance

type Instance struct {
	id InstanceID

	OnKeyDown func(Context, *KeyDown) error
}

func (i *Instance) ctx(ctx context.Context, conn *Conn) Context {
	return &instanceCtx{ctx, conn, i.id}
}
