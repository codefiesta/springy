package app

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.springy.io/model"
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

func Start() {

}

func Perform(client *Client, request *model.DatabaseRequest) {
	fmt.Printf("✅ Performing: %#v\n", request)

	switch request.Action {
	case model.Read:
		// Performs a single find request for the client
		_find(client, request)
		break
	case model.Write:
		switch request.Operation {
		case model.Insert:
			_insert(client, request)
			break
		case model.Update:
			_update(client, request)
			break
		case model.Delete:
			_delete(client, request)
			break
		case model.Replace:
			_replace(client, request)
			break
		}
		break
	case model.Watch:
		// Performs a change stream watch
		_watch(client, request)
		break
	}

}

func _find(client *Client, request *model.DatabaseRequest) {

}

func _insert(client *Client, request *model.DatabaseRequest) {
	log.Println("🚨 Inserting [", request.Collection, "]")
	collection := database.Collection(request.Collection)
	result, err := collection.InsertOne(context.Background(), request.Value)

	if err != nil {
		log.Fatal(err)
	}
	documentID := result.InsertedID
	var data = bson.M{
		"_sid": request.Identifier,
		"_id": documentID,
	}
	// Send the ID back
	fmt.Printf("⭐️ %v\n", documentID)
	go client.OnData(data)
}

func _update(client *Client, request *model.DatabaseRequest) {
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

func _delete(client *Client, request *model.DatabaseRequest) {

}

func _replace(client *Client, request *model.DatabaseRequest) {

}

// Starts watching (observing) a change stream
func _watch(client *Client, request *model.DatabaseRequest) {

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
		data["_sid"] = identifier
		//fmt.Printf("🍄 %v\n", data)
		client.OnData(data)
	}
}
