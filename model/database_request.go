package model

// Encapsulates our basic database request
type DatabaseRequest struct {
	Path      string                 `json:"path"`
	Action    DatabaseAction         `json:"action"`
	Operation DatabaseOperation      `json:"operation"`
	Value     map[string]interface{} `json:"value"`
}
