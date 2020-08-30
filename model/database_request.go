package model

// Encapsulates our basic database request
type DatabaseRequest struct {

	// The request identifier
	Identifier string `json:"_sid"`
	// The database collection name
	Collection string `json:"collection"`
	// The key of a document inside the collection (optional)
	Key string `json:"key"`
	// The action type
	Action DatabaseAction `json:"action"`
	/// The data operation
	Operation DatabaseOperation `json:"operation"`
	// The document value (optional)
	Value map[string]interface{} `json:"value"`
}
