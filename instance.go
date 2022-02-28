package streamdeck

import (
	"runtime/debug"
	"sync"

	"golang.org/x/net/context"
)

type InstanceContext interface {
	Context

	ShowOK() error
	ShowAlert() error
}

type instanceCtx struct {
	context.Context

	sdk        *SDK
	instanceID InstanceID
}

func (ctx *instanceCtx) OpenURL(url string) error {
	return ctx.sdk.OpenURL(url)
}

func (ctx *instanceCtx) ShowOK() error {
	return ctx.sdk.ShowOK(ctx.instanceID)
}

func (ctx *instanceCtx) ShowAlert() error {
	return ctx.sdk.ShowAlert(ctx.instanceID)
}

func (ctx *instanceCtx) Log(a ...interface{}) {
	ctx.sdk.Log(a...)
}

func (ctx *instanceCtx) Logf(format string, a ...interface{}) {
	ctx.sdk.Logf(format, a...)
}

type InstanceFactory func(ctx InstanceContext, id InstanceID) *Instance

type Instance struct {
	id  InstanceID
	sdk *SDK

	// Use slice + chan instead chan Event to have unlimited size of mailbox.
	mailbox []Event
	notify  chan struct{}
	mu      sync.Mutex

	OnKeyDown func(InstanceContext, *KeyDown) error
	OnKeyUp   func(InstanceContext, *KeyUp) error
}

func (i *Instance) ctx(ctx context.Context, sdk *SDK) InstanceContext {
	return &instanceCtx{ctx, sdk, i.id}
}

func (i *Instance) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-i.notify:
		}

		i.mu.Lock()
		events := append([]Event(nil), i.mailbox...)
		i.mailbox = i.mailbox[:0]
		i.mu.Unlock()

		for _, event := range events {
			i.sdk.debugf("[DEBUG] go-stream-deck-sdk: instance(%s) received: %T", i.id, event)

			var err error
			switch event := event.(type) {
			case *KeyDown:
				if i.OnKeyDown != nil {
					err = i.OnKeyDown(i.ctx(ctx, i.sdk), event)
				}
			case *KeyUp:
				if i.OnKeyUp != nil {
					err = i.OnKeyUp(i.ctx(ctx, i.sdk), event)
				}
			}
			if err != nil {
				i.sdk.Logf("go-stream-deck-sdk: error on instance(%s): %T: %v", i.id, event, err)
			}
		}
	}
}

func (i *Instance) handle(ev Event) {
	i.mu.Lock()
	i.mailbox = append(i.mailbox, ev)
	i.mu.Unlock()

	select {
	case i.notify <- struct{}{}:
	default:
	}
}

type supervisor struct {
	ctx          context.Context
	mu           sync.Mutex
	instanceByID map[InstanceID]*Instance
	sdk          *SDK
	factory      InstanceFactory
}

func newSupervisor(ctx context.Context, sdk *SDK, f InstanceFactory) *supervisor {
	return &supervisor{
		ctx,
		sync.Mutex{},
		make(map[InstanceID]*Instance),
		sdk,
		f,
	}
}

func (s *supervisor) handle(ev Event) {
	switch ev := ev.(type) {
	case *DidReceiveSettings:
		s.tell(ev.Context, ev)
	case *DidReceiveGlobalSettings:
		s.tellAll(ev)
	case *KeyDown:
		s.tell(ev.Context, ev)
	case *KeyUp:
		s.tell(ev.Context, ev)
	case *WillAppear:
		s.tell(ev.Context, ev)
	case *WillDisappear:
		s.tell(ev.Context, ev)
	case *TitleParametersDidChange:
		s.tell(ev.Context, ev)
	case *DeviceDidConnect:
		s.tellAll(ev)
	case *DeviceDidDisconnect:
		s.tellAll(ev)
	case *ApplicationDidLaunch:
		s.tellAll(ev)
	case *ApplicationDidTerminate:
		s.tellAll(ev)
	case *SystemDidWakeUp:
		s.tellAll(ev)
	case *PropertyInspectorDidAppear:
		s.tell(ev.Context, ev)
	case *PropertyInspectorDidDisappear:
		s.tell(ev.Context, ev)
	case *SendToPlugin:
		s.tell(ev.Context, ev)
	default:
		s.sdk.Logf("go-stream-deck-sdk: unknown event: %T", ev)
	}
}

func (s *supervisor) tell(id InstanceID, ev Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	instance, ok := s.instanceByID[id]
	if !ok {
		s.sdk.debugf("[DEBUG] go-stream-deck-sdk: spawn instance(%s): %T", id, ev)
		instance = s.spawn(id)
		if instance == nil {
			s.sdk.Logf("go-stream-deck-sdk: no instance returned: id = %s", id)
			return
		}
		s.instanceByID[id] = instance
	}

	instance.handle(ev)
}

func (s *supervisor) tellAll(ev Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, instance := range s.instanceByID {
		instance.handle(ev)
	}
}

func (s *supervisor) spawn(id InstanceID) *Instance {
	instance := s.factory(&instanceCtx{s.ctx, s.sdk, id}, id)
	if instance == nil {
		return nil
	}
	instance.id = id
	instance.notify = make(chan struct{}, 1)
	instance.sdk = s.sdk
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.sdk.Logf("go-stream-deck-sdk: panic on instance(%v): %v", instance.id, r)
				s.sdk.Log(string(debug.Stack()))
				_ = s.sdk.ShowAlert(instance.id)
				s.mu.Lock()
				delete(s.instanceByID, instance.id)
				s.mu.Unlock()
			}
		}()

		instance.run(s.ctx)
	}()

	return instance
}
