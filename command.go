package streamdeck

import (
	"encoding/json"
	"fmt"
)

// Command interface is to restrict command object that can
// implement the interface. The commands that are implementing
// the interface is enumerated below the definition in the source code.
type Command interface {
	commandMark()
	Event() string
}

var _ = []Command{
	(*LogMessage)(nil),
}

type commandMarkImpl struct{}

func (*commandMarkImpl) commandMark() {}

type commandPayload struct {
	Event   string          `json:"event"`
	Action  string          `json:"action,omitempty"`
	Context string          `json:"context,omitempty"`
	Device  string          `json:"device,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

func newCommandPayload(cmd Command) (*commandPayload, error) {
	p := &commandPayload{
		Event: cmd.Event(),
	}

	payload, err := json.Marshal(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal a command: %w: %v", err, cmd)
	}

	p.Payload = payload

	if t, ok := cmd.(interface{ GetAction() string }); ok {
		p.Action = t.GetAction()
	}
	if t, ok := cmd.(interface{ GetContext() string }); ok {
		p.Action = t.GetContext()
	}
	if t, ok := cmd.(interface{ GetDevice() string }); ok {
		p.Action = t.GetDevice()
	}

	return p, nil
}

type LogMessage struct {
	commandMarkImpl

	Message string `json:"message"`
}

func (l *LogMessage) Event() string {
	return "logMessage"
}
