package streamdeck

import (
	"encoding/base64"
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
	(*OpenURL)(nil),
	(*LogMessage)(nil),
	(*SetTitle)(nil),
	(*SetImage)(nil),
	(*ShowAlert)(nil),
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

type OpenURL struct {
	payloadCommand

	URL string `json:"url"`
}

func (*OpenURL) event() string {
	return "openUrl"
}

type LogMessage struct {
	payloadCommand

	Message string `json:"message"`
}

func (*LogMessage) event() string {
	return "logMessage"
}

type SetTitle struct {
	payloadCommand

	Context InstanceID `json:"-"`

	Title  string `json:"title,omitempty"`
	Target Target `json:"target"`
	State  int    `json:"state"`
}

func (*SetTitle) event() string {
	return "setTitle"
}

func (cmd *SetTitle) getContext() string {
	return string(cmd.Context)
}

type SetImage struct {
	payloadCommand

	Context InstanceID `json:"-"`

	Image  Image  `json:"image,omitempty"`
	Target Target `json:"target"`
	State  int    `json:"state"`
}

func (*SetImage) event() string {
	return "setImage"
}

func (cmd *SetImage) getContext() string {
	return string(cmd.Context)
}

type ShowAlert struct {
	noPayloadCommand

	Context InstanceID `json:"-"`
}

func (*ShowAlert) event() string {
	return "showAlert"
}

func (cmd *ShowAlert) getContext() string {
	return string(cmd.Context)
}

type ShowOK struct {
	noPayloadCommand

	Context InstanceID `json:"-"`
}

func (*ShowOK) event() string {
	return "showOk"
}

func (cmd *ShowOK) getContext() string {
	return string(cmd.Context)
}

type Target int

const (
	TargetBoth     Target = 0
	TargetHardware Target = 1
	TargetSoftware Target = 2
)

// Image represents Data URI scheme for image data.
type Image string

func NewImage(filetype string, data []byte) Image {
	return Image(fmt.Sprintf("data:image/%s;base64,%s", filetype, base64.StdEncoding.EncodeToString(data)))
}
