package app

import (
	"bytes"
	"encoding/json"
)

type Scope int

const (
	Find Scope = iota
	FindOne
	Write
	Watch
)

func (scope Scope) String() string {
	return scopeValue[scope]
}

var scopeValue = map[Scope]string{
	Find:    "find",
	FindOne: "findOne",
	Write:   "write",
	Watch:   "watch",
}

var scopeID = map[string]Scope{
	"find":    Find,
	"findOne": FindOne,
	"write":   Write,
	"watch":   Watch,
}

// MarshalJSON marshals the enum as a quoted json string
func (scope Scope) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(scopeValue[scope])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (scope *Scope) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*scope = scopeID[j]
	return nil
}
