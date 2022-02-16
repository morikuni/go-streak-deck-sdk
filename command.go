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
	event() string
}

var _ = []Command{
	(*LogMessage)(nil),
	(*ShowOK)(nil),
}

type noPayloadCommand struct{}

func (*noPayloadCommand) commandMark() {}

type payloadCommand struct{}

func (*payloadCommand) commandMark() {}

func (*payloadCommand) hasPayload() {}

type commandPayload struct {
	Event   string          `json:"event"`
	Context string          `json:"context,omitempty"`
	Action  string          `json:"action,omitempty"`
	Device  string          `json:"device,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

func newCommandPayload(cmd Command, pluginUUID string) (*commandPayload, error) {
	p := &commandPayload{
		Event:   cmd.event(),
		Context: pluginUUID,
	}

	if _, ok := cmd.(interface{ hasPayload() }); ok {
		payload, err := json.Marshal(cmd)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal a command: %w: %v", err, cmd)
		}
		p.Payload = payload
	}
	if t, ok := cmd.(interface{ getContext() string }); ok {
		p.Context = t.getContext()
	}
	if t, ok := cmd.(interface{ getAction() string }); ok {
		p.Action = t.getAction()
	}
	if t, ok := cmd.(interface{ getDevice() string }); ok {
		p.Device = t.getDevice()
	}

	return p, nil
}

type LogMessage struct {
	payloadCommand

	Message string `json:"message"`
}

func (*LogMessage) event() string {
	return "logMessage"
}

type ShowOK struct {
	noPayloadCommand

	Context InstanceID
}

func (*ShowOK) event() string {
	return "showOk"
}

func (cmd *ShowOK) getContext() string {
	return string(cmd.Context)
}
