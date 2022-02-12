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
}

type eventMarkImpl struct{}

func (*eventMarkImpl) eventMark() {}

type eventPayload struct {
	Action  string          `json:"action"`
	Event   string          `json:"event"`
	Context string          `json:"context"`
	Device  string          `json:"device"`
	Payload json.RawMessage `json:"payload"`
}

func (p eventPayload) Typed() (Event, error) {
	switch p.Event {
	case "keyDown":
		e := &KeyDown{
			Action:  p.Action,
			Context: p.Context,
			Device:  p.Device,
		}
		err := json.Unmarshal(p.Payload, e)
		if err != nil {
			return nil, fmt.Errorf("failed to bind event to %T: %w", e, err)
		}
		return e, nil
	default:
		// TODO: return error
		fmt.Println("unknown event", p)
		return &KeyDown{Action: p.Event + " " + p.Action, Context: p.Context, Device: p.Device}, nil
	}
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

type Coordinates struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}
