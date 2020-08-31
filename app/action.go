package app

import (
	"bytes"
	"encoding/json"
)

type Action int

const (
	Watch Action = iota
	Read
	Write
)

func (action Action) String() string {
	return actionValue[action]
}

var actionValue = map[Action]string{
	Watch: "watch",
	Read:  "read",
	Write: "write",
}

var actionID = map[string]Action{
	"watch": Watch,
	"read":  Read,
	"write": Write,
}

// MarshalJSON marshals the enum as a quoted json string
func (action Action) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(actionValue[action])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (action *Action) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*action = actionID[j]
	return nil
}
