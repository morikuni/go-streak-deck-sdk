package main

import (
	"encoding/json"
	"os"

	"github.com/morikuni/go-stream-deck-sdk/manifest"
)

func main() {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	err := enc.Encode(&manifest.Manifest{
		Actions: []manifest.Action{
			{
				Icon:                  "icon",
				Name:                  "Hello World",
				PropertyInspectorPath: nil,
				States: []manifest.State{
					{
						Image:            "icon",
						MultiActionImage: nil,
						Name:             nil,
						Title:            nil,
						ShowTitle:        nil,
						TitleColor:       nil,
						TitleAlignment:   nil,
						FontFamily:       nil,
						FontStyle:        nil,
						FontSize:         nil,
						FontUnderline:    nil,
					},
				},
				SupportedInMultiActions: nil,
				Tooltip:                 nil,
				UUID:                    "com.github.morikuni.helloworld",
				VisibleInActionsList:    nil,
			},
		},
		Author:                "morikuni",
		Category:              nil,
		CategoryIcon:          nil,
		CodePath:              "helloworld",
		CodePathMac:           nil,
		CodePathWin:           nil,
		Description:           "hello world app",
		Icon:                  "icon",
		Name:                  "Hello World",
		Profiles:              nil,
		PropertyInspectorPath: nil,
		DefaultWindowSize:     nil,
		URL:                   nil,
		Version:               "0.0.0",
		SDKVersion:            2,
		OS: []manifest.OS{
			{
				Platform:       manifest.PlatformMac,
				MinimumVersion: "10",
			},
		},
		Software: manifest.Software{
			MinimumVersion: "5.0",
		},
		ApplicationsToMonitor: nil,
	})
	if err != nil {
		panic(err)
	}
}
