package model

import (
	"bytes"
	"encoding/json"
)

type DatabaseAction int

const (
	Read DatabaseAction = iota
	Write
	Watch
)

func (action DatabaseAction) String() string {
	return databaseActionValue[action]
}

var databaseActionValue = map[DatabaseAction]string{
	Read:  "read",
	Write: "write",
	Watch: "watch",
}

var databaseActionID = map[string]DatabaseAction{
	"read":  Read,
	"write": Write,
	"watch": Watch,
}

// MarshalJSON marshals the enum as a quoted json string
func (action DatabaseAction) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(databaseActionValue[action])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (action *DatabaseAction) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*action = databaseActionID[j]
	return nil
}
