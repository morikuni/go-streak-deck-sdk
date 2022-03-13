package streamdeck

import (
	"runtime/debug"
	"sync"
)

type InstanceContext interface {
	Context

	SetTitle(title string, target Target, state int) error
	SetImage(img Image, target Target, state int) error
	ShowAlert() error
	ShowOK() error
}

type instanceCtx struct {
	Context
	instanceID InstanceID
}

func (ctx *instanceCtx) OpenURL(url string) error {
	return ctx.SDK().OpenURL(url)
}

func (ctx *instanceCtx) SetTitle(title string, target Target, state int) error {
	return ctx.SDK().SetTitle(ctx.instanceID, title, target, state)
}

func (ctx *instanceCtx) SetImage(img Image, target Target, state int) error {
	return ctx.SDK().SetImage(ctx.instanceID, img, target, state)
}

func (ctx *instanceCtx) ShowAlert() error {
	return ctx.SDK().ShowAlert(ctx.instanceID)
}

func (ctx *instanceCtx) ShowOK() error {
	return ctx.SDK().ShowOK(ctx.instanceID)
}

type InstanceFactory func(ctx InstanceContext, id InstanceID) *Instance

type Instance struct {
	id InstanceID

	// Use slice + chan instead chan Event to have unlimited size of mailbox.
	mailbox []Event
	notify  chan struct{}
	mu      sync.Mutex

	OnKeyDown func(InstanceContext, *KeyDown) error
	OnKeyUp   func(InstanceContext, *KeyUp) error
}

func (i *Instance) ctx(ctx Context) InstanceContext {
	return &instanceCtx{ctx, i.id}
}

func (i *Instance) run(ctx Context) {
	ictx := i.ctx(ctx)
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
			ictx.SDK().debugf("[DEBUG] go-stream-deck-sdk: instance(%s) received: %T", i.id, event)

			var err error
			switch event := event.(type) {
			case *KeyDown:
				if i.OnKeyDown != nil {
					err = i.OnKeyDown(ictx, event)
				}
			case *KeyUp:
				if i.OnKeyUp != nil {
					err = i.OnKeyUp(ictx, event)
				}
			}
			if err != nil {
				ictx.Logf("go-stream-deck-sdk: error on instance(%s): %T: %v", i.id, event, err)
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

type Mux struct {
	mu           sync.Mutex
	instanceByID map[InstanceID]*Instance
	factory      InstanceFactory
}

func NewMux(f InstanceFactory) *Mux {
	return &Mux{
		sync.Mutex{},
		make(map[InstanceID]*Instance),
		f,
	}
}

func (s *Mux) Handle(ctx Context, ev Event) error {
	switch ev := ev.(type) {
	case *DidReceiveSettings:
		s.tell(ctx, ev.Context, ev)
	case *DidReceiveGlobalSettings:
		s.tellAll(ev)
	case *KeyDown:
		s.tell(ctx, ev.Context, ev)
	case *KeyUp:
		s.tell(ctx, ev.Context, ev)
	case *WillAppear:
		s.tell(ctx, ev.Context, ev)
	case *WillDisappear:
		s.tell(ctx, ev.Context, ev)
	case *TitleParametersDidChange:
		s.tell(ctx, ev.Context, ev)
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
		s.tell(ctx, ev.Context, ev)
	case *PropertyInspectorDidDisappear:
		s.tell(ctx, ev.Context, ev)
	case *SendToPlugin:
		s.tell(ctx, ev.Context, ev)
	default:
		ctx.Logf("go-stream-deck-sdk: unknown event: %T", ev)
	}
	return nil
}

func (s *Mux) tell(ctx Context, id InstanceID, ev Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	instance, ok := s.instanceByID[id]
	if !ok {
		ctx.SDK().debugf("[DEBUG] go-stream-deck-sdk: spawn instance(%s): %T", id, ev)
		instance = s.spawn(ctx, id)
		if instance == nil {
			ctx.Logf("go-stream-deck-sdk: no instance returned: id = %s", id)
			return
		}
		s.instanceByID[id] = instance
	}

	instance.handle(ev)
}

func (s *Mux) tellAll(ev Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, instance := range s.instanceByID {
		instance.handle(ev)
	}
}

func (s *Mux) spawn(ctx Context, id InstanceID) *Instance {
	instance := s.factory(&instanceCtx{ctx, id}, id)
	if instance == nil {
		return nil
	}
	instance.id = id
	instance.notify = make(chan struct{}, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				ctx.Logf("go-stream-deck-sdk: panic on instance(%v): %v", instance.id, r)
				ctx.Log(string(debug.Stack()))
				_ = ctx.SDK().ShowAlert(instance.id)
				s.mu.Lock()
				delete(s.instanceByID, instance.id)
				s.mu.Unlock()
			}
		}()

		instance.run(ctx)
	}()

	return instance
}
