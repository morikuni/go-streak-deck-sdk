package streamdeck

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEventPayload_Typed(t *testing.T) {
	for name, tt := range map[string]struct {
		json string

		want Event
	}{
		"keyDown": {
			keyDownJSON,
			&KeyDown{
				Action:   "com.elgato.example.action1",
				Context:  "context",
				Device:   "device",
				Settings: json.RawMessage(`{}`),
				Coordinates: Coordinates{
					Row:    1,
					Column: 3,
				},
				State:            1,
				UserDesiredState: 1,
				IsInMultiAction:  true,
			},
		},
		"keyUp": {
			keyUpJSON,
			&KeyUp{
				Action:   "com.elgato.example.action1",
				Context:  "context",
				Device:   "device",
				Settings: json.RawMessage(`{}`),
				Coordinates: Coordinates{
					Row:    1,
					Column: 3,
				},
				State:            1,
				UserDesiredState: 1,
				IsInMultiAction:  true,
			},
		},
		"willAppear": {
			willAppearJSON,
			&WillAppear{
				Action:   "com.elgato.example.action1",
				Context:  "context",
				Device:   "device",
				Settings: json.RawMessage(`{}`),
				Coordinates: Coordinates{
					Row:    1,
					Column: 3,
				},
				State:           1,
				IsInMultiAction: true,
			},
		},
		"willDisappear": {
			willDisappearJSON,
			&WillDisappear{
				Action:   "com.elgato.example.action1",
				Context:  "context",
				Device:   "device",
				Settings: json.RawMessage(`{}`),
				Coordinates: Coordinates{
					Row:    1,
					Column: 3,
				},
				State:           1,
				IsInMultiAction: true,
			},
		},
		"titleParameterDidChange": {
			titleParameterDidChangeJSON,
			&TitleParametersDidChange{
				Action:   "com.elgato.example.action1",
				Context:  "context",
				Device:   "device",
				Settings: json.RawMessage(`{}`),
				Coordinates: Coordinates{
					Row:    1,
					Column: 3,
				},
				State: 1,
				Title: "title",
				TitleParameters: TitleParameters{
					FontFamily:     "fontFamily",
					FontSize:       12,
					FontStyle:      "fontStyle",
					FontUnderline:  true,
					ShowTitle:      true,
					TitleAlignment: AlignmentBottom,
					TitleColor:     "#ffffff",
				},
			},
		},
		"deviceDidConnect": {
			deviceDidConnectJSON,
			&DeviceDidConnect{
				Device: "device",
				DeviceInfo: &DeviceInfo{
					Name: "Device Name",
					Type: DeviceTypeStreamDeckMini,
					Size: Size{
						Rows:    3,
						Columns: 5,
					},
				},
			},
		},
		"deviceDidDisconnect": {
			deviceDidDisconnectJSON,
			&DeviceDidDisconnect{
				Device: "device",
			},
		},
		"applicationDidLaunch": {
			applicationDidLaunchJSON,
			&ApplicationDidLaunch{
				Application: "com.apple.mail",
			},
		},
		"applicationDidTerminate": {
			applicationDidTerminateJSON,
			&ApplicationDidTerminate{
				Application: "com.apple.mail",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			var ep eventPayload
			err := json.Unmarshal([]byte(tt.json), &ep)
			noError(t, err)

			tp, err := ep.Typed()
			noError(t, err)

			equal(t, tp, tt.want, ignoreUnexported(tt.want))
		})
	}
}

var keyDownJSON = `{
    "action": "com.elgato.example.action1",
    "event": "keyDown",
    "context": "context",
    "device": "device",
    "payload": {
        "settings": {},
        "coordinates": {
            "column": 3, 
            "row": 1
        },
        "state": 1,
        "userDesiredState": 1,
        "isInMultiAction": true
    }
}`

var keyUpJSON = `{
    "action": "com.elgato.example.action1",
    "event": "keyUp",
    "context": "context",
    "device": "device",
    "payload": {
        "settings": {},
        "coordinates": {
            "column": 3, 
            "row": 1
        },
        "state": 1,
        "userDesiredState": 1,
        "isInMultiAction": true
    }
}`

var willAppearJSON = `{
    "action": "com.elgato.example.action1",
    "event": "willAppear",
    "context": "context",
    "device": "device",
    "payload": {
        "settings": {},
        "coordinates": {
            "column": 3, 
            "row": 1
        },
        "state": 1,
        "isInMultiAction": true
    }
}`

var willDisappearJSON = `{
    "action": "com.elgato.example.action1",
    "event": "willDisappear",
    "context": "context",
    "device": "device",
    "payload": {
        "settings": {},
        "coordinates": {
            "column": 3, 
            "row": 1
        },
        "state": 1,
        "isInMultiAction": true
    }
}`

var titleParameterDidChangeJSON = `{
  "action": "com.elgato.example.action1", 
  "event": "titleParametersDidChange", 
  "context": "context", 
  "device": "device", 
  "payload": {
    "coordinates": {
      "column": 3, 
      "row": 1
    }, 
    "settings": {}, 
    "state": 1, 
    "title": "title", 
    "titleParameters": {
      "fontFamily": "fontFamily", 
      "fontSize": 12, 
      "fontStyle": "fontStyle", 
      "fontUnderline": true, 
      "showTitle": true, 
      "titleAlignment": "bottom", 
      "titleColor": "#ffffff"
    }
  }
}`

var deviceDidConnectJSON = `{
    "event": "deviceDidConnect",
    "device": "device",
    "deviceInfo": {
        "name": "Device Name",
        "type": 1,
        "size": {
            "rows": 3,
            "columns": 5
        }
    }
}`

var deviceDidDisconnectJSON = `{
    "event": "deviceDidDisconnect",
    "device": "device"
}`

var applicationDidLaunchJSON = `{
    "event": "applicationDidLaunch",
    "payload" : {
        "application": "com.apple.mail"
    }
}`

var applicationDidTerminateJSON = `{
    "event": "applicationDidTerminate",
    "payload" : {
        "application": "com.apple.mail"
    }
}`

func noError(tb testing.TB, err error) {
	tb.Helper()

	if err != nil {
		tb.Fatal("unexpected error:", err)
	}
}

func equal(tb testing.TB, got, want interface{}, opts ...cmp.Option) {
	tb.Helper()

	if diff := cmp.Diff(got, want, opts...); diff != "" {
		tb.Fatalf("(+want, -got): %s", diff)
	}
}

func ignoreUnexported(v interface{}) cmp.Option {
	return cmpopts.IgnoreUnexported(reflect.Indirect(reflect.ValueOf(v)).Interface())
}
