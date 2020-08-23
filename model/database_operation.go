package model

import (
	"bytes"
	"encoding/json"
)

type DatabaseOperation int

const (
	Insert DatabaseOperation = iota
	Update
	Delete
	Replace
)

func (operation DatabaseOperation) String() string {
	return databaseOperationValue[operation]
}

var databaseOperationValue = map[DatabaseOperation]string{
	Insert:  "insert",
	Update:  "update",
	Delete:  "delete",
	Replace: "replace",
}

var databaseOperationID = map[string]DatabaseOperation{
	"insert":  Insert,
	"update":  Update,
	"delete":  Delete,
	"replace": Replace,
}

// MarshalJSON marshals the enum as a quoted json string
func (operation DatabaseOperation) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(databaseOperationValue[operation])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshalls a quoted json string to the enum value
func (operation *DatabaseOperation) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*operation = databaseOperationID[j]
	return nil
}
