package api

import (
	"bytes"
	"encoding/json"
)

type DocumentScope int

const (
	// Single Read Request
	Find DocumentScope = iota
	// Single Read Request
	FindOne
	// Single Write Request
	Write
	// Subscribe Request
	Watch
)

func (scope DocumentScope) String() string {
	return scopeValue[scope]
}

var scopeValue = map[DocumentScope]string{
	Find:    "find",
	FindOne: "findOne",
	Write:   "write",
	Watch:   "watch",
}

var scopeID = map[string]DocumentScope{
	"find":    Find,
	"findOne": FindOne,
	"write":   Write,
	"watch":   Watch,
}

// MarshalJSON marshals the enum as a quoted json string
func (scope DocumentScope) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(scopeValue[scope])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (scope *DocumentScope) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*scope = scopeID[j]
	return nil
}
