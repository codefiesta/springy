package app

// Encapsulates a basic request sent from a client
type Request struct {

	// The request unique identifier
	Uid string `json:"_uid"`
	// The database collection name
	Collection string `json:"collection"`
	// The key of a document inside the collection (optional)
	Key string `json:"key"`
	// The action to perform
	Action Action `json:"action"`
	// The operation to observe or perform
	Operation Operation `json:"operation"`
	// The document value (optional)
	Value map[string]interface{} `json:"value"`
	// Flag indicating if request should be processed on disconnect
	OnDisconnect bool `json:"onDisconnect"`
}