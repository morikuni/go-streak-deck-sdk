package streamdeck

import (
	"encoding/json"
	"fmt"
	"testing"
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
			equal(t, cp, tt.want)
		})
	}
}
