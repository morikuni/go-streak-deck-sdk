package streamdeck

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewCommandPayload(t *testing.T) {
	toJSON := func(i interface{}) json.RawMessage {
		b, err := json.Marshal(i)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}

	for _, tt := range []struct {
		cmd Command

		want *commandPayload
	}{
		{
			cmd: &OpenURL{
				URL: "url",
			},
			want: &commandPayload{
				Event:   "openUrl",
				Context: "pluginUUID", // don't need to set this field, but set automatically.
				Payload: toJSON(map[string]string{
					"url": "url",
				}),
			},
		},
		{
			cmd: &LogMessage{
				Message: "message",
			},
			want: &commandPayload{
				Event:   "logMessage",
				Context: "pluginUUID", // don't need to set this field, but set automatically.
				Payload: toJSON(map[string]string{
					"message": "message",
				}),
			},
		},
		{
			cmd: &SetTitle{
				Context: "instanceID",
				Title:   "title",
				Target:  TargetHardware,
				State:   2,
			},
			want: &commandPayload{
				Event:   "setTitle",
				Context: "instanceID",
				Payload: toJSON(map[string]interface{}{
					"title":  "title",
					"target": 1,
					"state":  2,
				}),
			},
		},
		{
			cmd: &SetImage{
				Context: "instanceID",
				Image:   NewImage("png", []byte("data")),
				Target:  TargetSoftware,
				State:   3,
			},
			want: &commandPayload{
				Event:   "setImage",
				Context: "instanceID",
				Payload: toJSON(map[string]interface{}{
					"image":  "data:image/png;base64,ZGF0YQ==",
					"target": 2,
					"state":  3,
				}),
			},
		},
		{
			cmd: &ShowAlert{
				Context: "instanceID",
			},
			want: &commandPayload{
				Event:   "showAlert",
				Context: "instanceID",
			},
		},
		{
			cmd: &ShowOK{
				Context: "instanceID",
			},
			want: &commandPayload{
				Event:   "showOk",
				Context: "instanceID",
			},
		},
	} {
		t.Run(fmt.Sprintf("%T", tt.cmd), func(t *testing.T) {
			cp, err := newCommandPayload(tt.cmd, "pluginUUID")
			noError(t, err)
			equal(t, cp, tt.want, cmpopts.IgnoreFields(commandPayload{}, "Payload"))
			equalJSON(t, cp.Payload, tt.want.Payload)
		})
	}
}
