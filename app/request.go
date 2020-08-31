package app

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Encapsulates a basic request sent from a client
type Request struct {

	// The request unique identifier
	Uid string `json:"_uid"`

	// The database collection name
	Collection string `json:"collection"`

	// The key of a document inside the collection (optional)
	Query map[string]interface{} `json:"query"`

	// The scope of work to perform
	Scope Scope `json:"scope"`

	// The operation to observe or perform
	Operation Operation `json:"operation"`

	// The document value (optional)
	Value map[string]interface{} `json:"value"`

	// Flag indicating if request should be processed on disconnect
	OnDisconnect bool `json:"onDisconnect"`
}

// Builds a document filter based on the query passed into the request
func (request *Request) filter() bson.M {

	var filters = bson.M{}
	for k, v := range request.Query {
		switch k {
		case "_id":
			docID, _ := primitive.ObjectIDFromHex(v.(string))
			filters[k] = docID
		default:
			filters[k] = v
		}
	}


	return filters
}