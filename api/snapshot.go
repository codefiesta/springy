package api

// Encapsulates a snapshot of a document in time
type DocumentSnapshot struct {

	// The document value (optional)
	Value map[string]interface{}
}
