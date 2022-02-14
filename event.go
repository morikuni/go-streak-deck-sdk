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
	(*DidReceiveSettings)(nil),
	(*DidReceiveGlobalSettings)(nil),
	(*KeyDown)(nil),
	(*KeyUp)(nil),
	(*WillAppear)(nil),
	(*WillDisappear)(nil),
	(*TitleParametersDidChange)(nil),
	(*DeviceDidConnect)(nil),
	(*DeviceDidDisconnect)(nil),
	(*ApplicationDidLaunch)(nil),
	(*ApplicationDidTerminate)(nil),
	(*SystemDidWakeUp)(nil),
	(*PropertyInspectorDidAppear)(nil),
	(*PropertyInspectorDidDisappear)(nil),
	(*SendToPlugin)(nil),
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
		case "didReceiveSettings":
			return &DidReceiveSettings{}
		case "didReceiveGlobalSettings":
			return &DidReceiveGlobalSettings{}
		case "keyDown":
			return &KeyDown{}
		case "keyUp":
			return &KeyUp{}
		case "willAppear":
			return &WillAppear{}
		case "willDisappear":
			return &WillDisappear{}
		case "titleParametersDidChange":
			return &TitleParametersDidChange{}
		case "deviceDidConnect":
			return &DeviceDidConnect{}
		case "deviceDidDisconnect":
			return &DeviceDidDisconnect{}
		case "applicationDidLaunch":
			return &ApplicationDidLaunch{}
		case "applicationDidTerminate":
			return &ApplicationDidTerminate{}
		case "systemDidWakeUp":
			return &SystemDidWakeUp{}
		case "propertyInspectorDidAppear":
			return &PropertyInspectorDidAppear{}
		case "propertyInspectorDidDisappear":
			return &PropertyInspectorDidDisappear{}
		case "sendToPlugin":
			return &SendToPlugin{}
		default:
			return nil
		}
	}()
	if e == nil {
		return nil, fmt.Errorf("unknown event: %s", ep.Event)
	}

	if len(ep.Payload) > 0 {
		err := json.Unmarshal(ep.Payload, e)
		if err != nil {
			return nil, fmt.Errorf("failed to bind event payload to %T: %w", e, err)
		}
	}

	err := json.Unmarshal(ep.Raw, e)
	if err != nil {
		return nil, fmt.Errorf("failed to bind event to %T: %w", e, err)
	}

	return e, nil
}

type DidReceiveSettings struct {
	eventMarkImpl

	Action  string `json:"action"`
	Context string `json:"context"`
	Device  string `json:"device"`

	Settings        json.RawMessage `json:"settings"`
	Coordinates     Coordinates     `json:"coordinates"`
	State           int             `json:"state"`
	IsInMultiAction bool            `json:"isInMultiAction"`
}

type DidReceiveGlobalSettings struct {
	eventMarkImpl

	Payload json.RawMessage `json:"payload"`
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

type TitleParametersDidChange struct {
	eventMarkImpl

	Action  string `json:"action"`
	Context string `json:"context"`
	Device  string `json:"device"`

	Settings        json.RawMessage `json:"settings"`
	Coordinates     Coordinates     `json:"coordinates"`
	State           int             `json:"state"`
	Title           string          `json:"title"`
	TitleParameters TitleParameters `json:"titleParameters"`
}

type DeviceDidConnect struct {
	eventMarkImpl

	Device     string      `json:"device"`
	DeviceInfo *DeviceInfo `json:"deviceInfo"`
}

type DeviceDidDisconnect struct {
	eventMarkImpl

	Device string `json:"device"`
}

type ApplicationDidLaunch struct {
	eventMarkImpl

	Application string `json:"application"`
}

type ApplicationDidTerminate struct {
	eventMarkImpl

	Application string `json:"application"`
}

type SystemDidWakeUp struct {
	eventMarkImpl
}

type PropertyInspectorDidAppear struct {
	eventMarkImpl

	Action  string `json:"action"`
	Context string `json:"context"`
	Device  string `json:"device"`
}

type PropertyInspectorDidDisappear struct {
	eventMarkImpl

	Action  string `json:"action"`
	Context string `json:"context"`
	Device  string `json:"device"`
}

type SendToPlugin struct {
	eventMarkImpl

	Action  string `json:"action"`
	Context string `json:"context"`

	Payload json.RawMessage `json:"payload"`
}

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

type TitleParameters struct {
	FontFamily     string    `json:"fontFamily"`
	FontSize       int       `json:"fontSize"`
	FontStyle      string    `json:"fontStyle"`
	FontUnderline  bool      `json:"fontUnderline"`
	ShowTitle      bool      `json:"showTitle"`
	TitleAlignment Alignment `json:"titleAlignment"`
	TitleColor     string    `json:"titleColor"`
}

type Alignment string

const (
	AlignmentTop    Alignment = "top"
	AlignmentBottom Alignment = "bottom"
	AlignmentMiddle Alignment = "middle"
)
