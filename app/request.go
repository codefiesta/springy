package app

import (
	"bytes"
	"encoding/json"
)

// Encapsulates a basic request sent from a client
type Request struct {

	// The request identifier
	Identifier string `json:"_uid"`
	// The database collection name
	Collection string `json:"collection"`
	// The key of a document inside the collection (optional)
	Key string `json:"key"`
	// The action type
	Action RequestAction `json:"action"`
	/// The data operation
	Operation DataOperation `json:"operation"`
	// The document value (optional)
	Value map[string]interface{} `json:"value"`
}

type RequestAction int

const (
	Read RequestAction = iota
	Write
	Watch
)

func (action RequestAction) String() string {
	return actionValue[action]
}

var actionValue = map[RequestAction]string{
	Read:  "read",
	Write: "write",
	Watch: "watch",
}

var actionID = map[string]RequestAction{
	"read":  Read,
	"write": Write,
	"watch": Watch,
}

// MarshalJSON marshals the enum as a quoted json string
func (action RequestAction) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(actionValue[action])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (action *RequestAction) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*action = actionID[j]
	return nil
}

type DataOperation int

const (
	Insert DataOperation = iota
	Update
	Delete
	Replace
)

func (operation DataOperation) String() string {
	return operationValue[operation]
}

var operationValue = map[DataOperation]string{
	Insert:  "insert",
	Update:  "update",
	Delete:  "delete",
	Replace: "replace",
}

var operationID = map[string]DataOperation{
	"insert":  Insert,
	"update":  Update,
	"delete":  Delete,
	"replace": Replace,
}

// MarshalJSON marshals the enum as a quoted json string
func (operation DataOperation) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(operationValue[operation])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshalls a quoted json string to the enum value
func (operation *DataOperation) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Insert' in this case.
	*operation = operationID[j]
	return nil
}
