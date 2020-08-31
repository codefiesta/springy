package app

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.springy.io/util"
	"log"
	"time"
)

var (
	database *mongo.Database
	config   *util.Configuration
)

func init() {

	fmt.Println("Connecting ...")
	config = util.Config()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Database.Uri))
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	database = client.Database(config.Database.Name)
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
}

func processRequest(client *Client, request *Request) {

	switch request.Action {
	case Read:
		// Performs a single find request
		_find(client, request)
		break
	case Write:
		// Performs a single CRUD operation
		switch request.Operation {
		case Insert:
			_insert(client, request)
			break
		case Update:
			_update(client, request)
			break
		case Delete:
			_delete(client, request)
			break
		case Replace:
			_replace(client, request)
			break
		}
		break
	case Watch:
		// Performs a change stream watch
		_watch(client, request)
		break
	}
}

func _find(client *Client, request *Request) {

	context := context.Background()
	collection := database.Collection(request.Collection)
	cursor, err := collection.Find(context, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = cursor.All(context, &results); err != nil {
		log.Fatal(err)
	}

	var snapshot = bson.M{
		"_uid":  request.Identifier,
		"key":   request.Collection,
		"value": results,
	}
	go client.writeResponse(snapshot)
}

func _insert(client *Client, request *Request) {
	collection := database.Collection(request.Collection)
	result, err := collection.InsertOne(context.Background(), request.Value)

	if err != nil {
		log.Fatal(err)
	}

	var snapshot = bson.M{
		"_uid": request.Identifier,
		"key":  result.InsertedID,
	}
	go client.writeResponse(snapshot)
}

func _update(client *Client, request *Request) {
	//collection := database.Collection(request.Collection)
	//id, _ := primitive.ObjectIDFromHex("5d9e0173c1305d2a54eb431a")
	//result, err := collection.UpdateOne(
	//	context.Background(),
	//	bson.M{"_id": id},
	//	bson.D{ bson.E{request.Value}},
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func _delete(client *Client, request *Request) {
	collection := database.Collection(request.Collection)
	docID, _ := primitive.ObjectIDFromHex(request.Key)
	var match = bson.M{"_id": docID}

	_, err := collection.DeleteOne(context.Background(), match)

	if err != nil {
		log.Fatal(err)
	}

	var snapshot = bson.M{
		"_uid": request.Identifier,
		"key":  request.Key,
	}
	go client.writeResponse(snapshot)
}

func _replace(client *Client, request *Request) {

}

// Starts watching (observing) a change stream
func _watch(client *Client, request *Request) {

	var matchingPipeline = bson.D{
		{
			"$match", bson.D{
			{"operationType", request.Operation.String()},
		},
		},
	}
	collection := database.Collection(request.Collection)
	collectionStream, err := collection.Watch(context.TODO(), mongo.Pipeline{matchingPipeline})

	if err != nil {
		log.Fatal(err)
	}

	streamContext, _ := context.WithCancel(context.Background())
	go _watchChangeStream(request.Identifier, client, streamContext, collectionStream)
}

func _watchChangeStream(identifier string, client *Client, context context.Context, stream *mongo.ChangeStream) {
	defer stream.Close(context)
	for stream.Next(context) {
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			panic(err)
		}

		docKey, _ := data["documentKey"].(bson.M)
		doc, _ := data["fullDocument"].(bson.M)
		var snapshot = bson.M{
			"_uid":  identifier,
			"key":   docKey["_id"],
			"value": doc,
		}
		client.writeResponse(snapshot)
	}
}
