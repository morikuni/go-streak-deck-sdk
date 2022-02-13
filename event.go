package streamdeck

import (
	"encoding/json"
	"fmt"
)

// Event interface is to restrict event object that can
// implement the interface. The events that are implementing
// the interface is enumerated below the definition in the source code.
type Event interface {
	eventMark()
}

var _ = []Event{
	(*KeyDown)(nil),
	(*DeviceDidConnect)(nil),
	(*DebugEvent)(nil),
}

type eventMarkImpl struct{}

func (*eventMarkImpl) eventMark() {}

type eventPayload struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
	Raw     json.RawMessage `json:"-"`
}

func (ep *eventPayload) UnmarshalJSON(bs []byte) error {
	var extractor struct {
		Event   string          `json:"event"`
		Payload json.RawMessage `json:"payload"`
	}

	err := json.Unmarshal(bs, &extractor)
	if err != nil {
		return err
	}

	ep.Event = extractor.Event
	ep.Payload = extractor.Payload
	ep.Raw = bs
	return nil
}

func (ep eventPayload) Typed() (Event, error) {
	var e = func() Event {
		switch ep.Event {
		case "keyDown":
			return &KeyDown{}
		case "keyUp":
			return &KeyUp{}
		case "willAppear":
			return &WillAppear{}
		case "willDisappear":
			return &WillDisappear{}
		case "deviceDidConnect":
			return &DeviceDidConnect{}
		default:
			return &DebugEvent{}
		}
	}()
	if e == nil {
		return nil, fmt.Errorf("unknown event: %s", ep.Event)
	}

	err := json.Unmarshal(ep.Raw, e)
	if err != nil {
		return nil, fmt.Errorf("failed to bind event to %T: %w", e, err)
	}

	if len(ep.Payload) > 0 {
		err = json.Unmarshal(ep.Payload, e)
		if err != nil {
			return nil, fmt.Errorf("failed to bind event payload to %T: %w", e, err)
		}
	}

	return e, nil
}

type KeyDown struct {
	eventMarkImpl

	Action  string `json:"action"`
	Context string `json:"context"`
	Device  string `json:"device"`

	Settings         json.RawMessage `json:"settings"`
	Coordinates      Coordinates     `json:"coordinates"`
	State            int             `json:"state"`
	UserDesiredState int             `json:"userDesiredState"`
	IsInMultiAction  bool            `json:"isInMultiAction"`
}

type KeyUp struct {
	eventMarkImpl

	Action  string `json:"action"`
	Context string `json:"context"`
	Device  string `json:"device"`

	Settings         json.RawMessage `json:"settings"`
	Coordinates      Coordinates     `json:"coordinates"`
	State            int             `json:"state"`
	UserDesiredState int             `json:"userDesiredState"`
	IsInMultiAction  bool            `json:"isInMultiAction"`
}

type WillAppear struct {
	eventMarkImpl

	Action  string `json:"action"`
	Context string `json:"context"`
	Device  string `json:"device"`

	Settings        json.RawMessage `json:"settings"`
	Coordinates     Coordinates     `json:"coordinates"`
	State           int             `json:"state"`
	IsInMultiAction bool            `json:"isInMultiAction"`
}

type WillDisappear struct {
	eventMarkImpl

	Action  string `json:"action"`
	Context string `json:"context"`
	Device  string `json:"device"`

	Settings        json.RawMessage `json:"settings"`
	Coordinates     Coordinates     `json:"coordinates"`
	State           int             `json:"state"`
	IsInMultiAction bool            `json:"isInMultiAction"`
}

type DeviceDidConnect struct {
	eventMarkImpl

	Device     string      `json:"device"`
	DeviceInfo *DeviceInfo `json:"deviceInfo"`
}

func (e *DeviceDidConnect) setDevice(s string) { e.Device = s }

type DeviceInfo struct {
	Name string     `json:"name"`
	Type DeviceType `json:"type"`
	Size Size       `json:"size"`
}

type DeviceType int

const (
	DeviceTypeStreamDeck       DeviceType = 0
	DeviceTypeStreamDeckMini   DeviceType = 1
	DeviceTypeStreamDeckXL     DeviceType = 2
	DeviceTypeStreamDeckMobile DeviceType = 3
	DeviceTypeCorsairGKeys     DeviceType = 4
	DeviceTypeStreamDeckPanel  DeviceType = 5
)

type Size struct {
	Rows    int `json:"rows"`
	Columns int `json:"columns"`
}

type Coordinates struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

// TODO: remove me
type DebugEvent struct {
	eventMarkImpl

	Event string `json:"event"`
	Raw   string `json:"-"`
}

func (e *DebugEvent) UnmarshalJSON(bs []byte) error {
	var eventExtractor struct {
		Event string `json:"event"`
	}

	err := json.Unmarshal(bs, &eventExtractor)
	if err != nil {
		return err
	}

	e.Event = eventExtractor.Event
	e.Raw = string(bs)
	return nil
}
