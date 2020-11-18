package app

import (
	"bytes"
	"encoding/json"
)

type Operation int

const (
	Insert Operation = iota
	Update
	Delete
	Replace
)

func (operation Operation) String() string {
	return operationValue[operation]
}

var operationValue = map[Operation]string{
	Insert:  "insert",
	Update:  "update",
	Delete:  "delete",
	Replace: "replace",
}

var operationID = map[string]Operation{
	"insert":  Insert,
	"update":  Update,
	"delete":  Delete,
	"replace": Replace,
}

// MarshalJSON marshals the enum as a quoted json string
func (operation Operation) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(operationValue[operation])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshalls a quoted json string to the enum value
func (operation *Operation) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Insert' in this case.
	*operation = operationID[j]
	return nil
}
