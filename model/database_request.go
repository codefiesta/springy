package model

// Encapsulates our basic database request
type DatabaseRequest struct {
	Path       string                 `json:"path"`
	Collection string                 `json:"collection"`
	Action     DatabaseAction         `json:"action"`
	Operation  DatabaseOperation      `json:"operation"`
	Value      map[string]interface{} `json:"value"`
}
