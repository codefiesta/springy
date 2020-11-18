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

	switch request.Scope {
	case Find:
		_find(client, request)
		break
	case FindOne:
		_findOne(client, request)
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

func _findOne(client *Client, request *Request) {
	context := context.Background()
	collection := database.Collection(request.Collection)
	result := collection.FindOne(context, request.filter())

	doc := bson.M{}
	if result.Err() == mongo.ErrNoDocuments {
		doc = bson.M{
			"value": bson.M{
				"_id": primitive.NewObjectID(),
			},
		}
	} else {
		if err := result.Decode(&doc); err != nil {
			panic(err)
		}
	}

	var response = bson.M{
		"_uid":       request.Uid,
		"_operation": request.Operation,
		"value":      doc,
	}
	go client.writeResponse(response)
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

	if request.OnDisconnect {
		return
	}

	var response = bson.M{
		"_uid":       request.Uid,
		"_operation": request.Operation,
		"value":      results,
	}
	go client.writeResponse(response)
}

func _insert(client *Client, request *Request) {
	collection := database.Collection(request.Collection)
	result, err := collection.InsertOne(context.Background(), request.Value)

	if err != nil {
		log.Fatal(err)
	}

	if request.OnDisconnect {
		return
	}

	request.Value["_id"] = result.InsertedID
	var response = bson.M{
		"_uid":       request.Uid,
		"_operation": request.Operation,
		"value":      request.Value,
	}
	go client.writeResponse(response)
}

func _update(client *Client, request *Request) {
	collection := database.Collection(request.Collection)
	result, err := collection.UpdateOne(context.Background(), request.filter(), request.Value)
	if err != nil {
		log.Fatal(err)
	}
	if request.OnDisconnect {
		return
	}
	request.Value["_id"] = result.UpsertedID
	var response = bson.M{
		"_uid":       request.Uid,
		"_operation": request.Operation,
		"value":      request.Value,
	}
	go client.writeResponse(response)
}

func _delete(client *Client, request *Request) {
	collection := database.Collection(request.Collection)
	_, err := collection.DeleteOne(context.Background(), request.filter())

	if err != nil {
		log.Fatal(err)
	}

	if request.OnDisconnect {
		return
	}

	var response = bson.M{
		"_uid":       request.Uid,
		"_operation": request.Operation,
		"value":      request.Query,
	}

	go client.writeResponse(response)
}

func _replace(client *Client, request *Request) {
	collection := database.Collection(request.Collection)
	result, err := collection.ReplaceOne(context.Background(), request.filter(), request.Value)
	if err != nil {
		log.Fatal(err)
	}
	if request.OnDisconnect {
		return
	}
	request.Value["_id"] = result.UpsertedID
	var response = bson.M{
		"_uid":       request.Uid,
		"_operation": request.Operation,
		"value":      request.Value,
	}
	go client.writeResponse(response)
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
	go _watchChangeStream(client, request, streamContext, collectionStream)
}

func _watchChangeStream(client *Client, request *Request, context context.Context, stream *mongo.ChangeStream) {
	defer stream.Close(context)
	for stream.Next(context) {
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			panic(err)
		}

		key, _ := data["documentKey"].(bson.M)
		doc, _ := data["fullDocument"].(bson.M)

		if doc == nil {
			switch request.Operation {
			case Update, Replace:
				request.Query["_id"] = key["_id"].(primitive.ObjectID).Hex()
				_findOne(client, request)
				return
			case Delete:
				doc = bson.M{
					"_id": key["_id"],
				}
				break
			default:
				break
			}
		}

		var response = bson.M{
			"_uid":       request.Uid,
			"_operation": request.Operation,
			"value":      doc,
		}
		go client.writeResponse(response)
	}
}
