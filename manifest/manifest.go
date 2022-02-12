// Package manifest defines type of manifest.json.
// The detailed document is here: https://developer.elgato.com/documentation/stream-deck/sdk/manifest/.
package manifest

import (
	"encoding/json"
)

type Manifest struct {
	Actions               []Action
	Author                string
	Category              *string `json:",omitempty"`
	CategoryIcon          *string `json:",omitempty"`
	CodePath              string
	CodePathMac           *string `json:",omitempty"`
	CodePathWin           *string `json:",omitempty"`
	Description           string
	Icon                  string
	Name                  string
	Profiles              []Profile `json:",omitempty"`
	PropertyInspectorPath *string   `json:",omitempty"`
	DefaultWindowSize     []int     `json:",omitempty"`
	URL                   *string   `json:",omitempty"`
	Version               string
	SDKVersion            int
	OS                    []OS
	Software              Software
	ApplicationsToMonitor *ApplicationsToMonitor `json:",omitempty"`
}

type Action struct {
	Icon                    string
	Name                    string
	PropertyInspectorPath   *string `json:",omitempty"`
	States                  []State
	SupportedInMultiActions *bool   `json:",omitempty"`
	Tooltip                 *string `json:",omitempty"`
	UUID                    string
	VisibleInActionsList    *bool `json:",omitempty"`
}

type State struct {
	Image            string
	MultiActionImage *string `json:",omitempty"`
	Name             *string `json:",omitempty"`
	Title            *string `json:",omitempty"`
	ShowTitle        *bool   `json:",omitempty"`
	TitleColor       *string `json:",omitempty"`
	TitleAlignment   *string `json:",omitempty"`
	FontFamily       *string `json:",omitempty"`
	FontStyle        *string `json:",omitempty"`
	FontSize         *int    `json:",omitempty"`
	FontUnderline    *bool   `json:",omitempty"`
}

type Profile struct {
	Name                        string
	DeviceType                  DeviceType
	ReadOnly                    *bool `json:",omitempty"`
	DontAutoSwitchWhenInstalled *bool `json:",omitempty"`
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

type OS struct {
	Platform       Platform
	MinimumVersion string
}

type Platform string

const (
	PlatformMac     Platform = "mac"
	PlatformWindows Platform = "windows"
)

type Software struct {
	MinimumVersion string
}

type ApplicationsToMonitor struct {
	Mac     []string
	Windows []string
}

func (a ApplicationsToMonitor) MarshalJSON() ([]byte, error) {
	type obj struct {
		Mac     []string `json:"mac"`
		Windows []string `json:"windows"`
	}
	o := obj(a)
	return json.Marshal(o)
}

func OptionalString(s string) *string {
	return &s
}

func OptionalBool(b bool) *bool {
	return &b
}

func OptionalInt(i int) *int {
	return &i
}
