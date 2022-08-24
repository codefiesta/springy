package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.springy.io/api"
	"go.springy.io/pkg/events"
	"go.springy.io/pkg/util"
	"log"
	"time"
)

var (
	database *mongo.Database
	env      *util.Environment
)

func init() {
	log.Println("🌱 [Initializing MongoDB] 🌱")
	env = util.Env()

	// https://github.com/mongodb/mongo-go-driver/blob/master/mongo/client_examples_test.go
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		AuthSource: env.Database.Db,
		Username:   env.Database.Username,
		Password:   env.Database.Password,
	}

	uri := env.Database.GetURI()

	clientOptions := options.Client().
		SetHosts([]string{uri}).
		SetDirect(true).
		SetAppName(env.Database.Db).
		SetAuth(credential).
		SetReplicaSet(env.Database.ReplicaSet).
		SetReadPreference(readpref.Primary())

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("💩 [Unable to ping mongo]: ", err)
	}

	database = client.Database(env.Database.Db)
	databases, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal("💩 [Unable to list mongo databases]: ", err)
	}
	log.Println("🌱", databases, "🌱")
}

func Run() {
	subscriber := make(chan events.Event)
	events.Subscribe(events.Mongo, subscriber)
	for {
		select {
		case e := <-subscriber:
			go handle(e)
		}
	}
}

// Processes a document request event
func handle(e events.Event) {
	// Make sure we are dealing with an API request
	if request, ok := e.Data.(api.DocumentRequest); ok {
		switch request.Scope {
		case api.Find:
			_find(e.Sender, request)
			break
		case api.FindOne:
			_findOne(e.Sender, request)
		case api.Write:
			// Performs a single CRUD operation
			switch request.Operation {
			case api.Insert:
				_insert(e.Sender, request)
				break
			case api.Update:
				_update(e.Sender, request)
				break
			case api.Delete:
				_delete(e.Sender, request)
				break
			case api.Replace:
				_replace(e.Sender, request)
				break
			}
			break
		case api.Watch:
			// Performs a change stream watch
			_watch(e.Sender, request)
			break
		}
	}
}

func publish(sender interface{}, doc bson.M) {
	snapshot := api.DocumentSnapshot{
		Value: doc,
	}
	go events.Publish(events.Websocket, sender, snapshot)
}

func _findOne(sender interface{}, request api.DocumentRequest) {
	context := context.Background()
	collection := database.Collection(request.Collection)
	result := collection.FindOne(context, request.Filter())

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

	snapshot := bson.M{
		"_uid":       request.Uid,
		"_operation": request.Operation,
		"value":      doc,
	}
	publish(sender, snapshot)
}

func _find(sender interface{}, request api.DocumentRequest) {

	context := context.Background()
	collection := database.Collection(request.Collection)
	cursor, err := collection.Find(context, bson.M{})
	if err != nil {
		log.Fatal("💩 [Unable to open collection]: ", err)
	}
	var results []bson.M
	if err = cursor.All(context, &results); err != nil {
		log.Fatal("💩 [Unable to open cursor]: ", err)
	}

	if request.OnDisconnect {
		return
	}

	snapshot := bson.M{
		"_uid":       request.Uid,
		"_operation": request.Operation,
		"value":      results,
	}
	publish(sender, snapshot)
}

func _insert(sender interface{}, request api.DocumentRequest) {
	collection := database.Collection(request.Collection)
	result, err := collection.InsertOne(context.Background(), request.Value)

	if err != nil {
		log.Fatal("💩 [Unable to insert]: ", err)
	}

	if request.OnDisconnect {
		return
	}

	request.Value["_id"] = result.InsertedID

	snapshot := bson.M{
		"_uid":       request.Uid,
		"_operation": request.Operation,
		"value":      request.Value,
	}
	publish(sender, snapshot)
}

func _update(sender interface{}, request api.DocumentRequest) {
	collection := database.Collection(request.Collection)
	result, err := collection.UpdateOne(context.Background(), request.Filter(), request.Value)
	if err != nil {
		log.Fatal("💩 [Unable to update]: ", err)
	}
	if request.OnDisconnect {
		return
	}
	request.Value["_id"] = result.UpsertedID

	snapshot := bson.M{
		"_uid":       request.Uid,
		"_operation": request.Operation,
		"value":      request.Value,
	}
	publish(sender, snapshot)
}

func _delete(sender interface{}, request api.DocumentRequest) {
	collection := database.Collection(request.Collection)
	_, err := collection.DeleteOne(context.Background(), request.Filter())

	if err != nil {
		log.Fatal("💩 [Unable to delete]: ", err)
	}

	if request.OnDisconnect {
		return
	}

	snapshot := bson.M{
		"_uid":       request.Uid,
		"_operation": request.Operation,
		"value":      request.Query,
	}

	publish(sender, snapshot)
}

func _replace(sender interface{}, request api.DocumentRequest) {
	collection := database.Collection(request.Collection)
	result, err := collection.ReplaceOne(context.Background(), request.Filter(), request.Value)
	if err != nil {
		log.Fatal("💩 [Unable to replace]: ", err)
	}

	if request.OnDisconnect {
		return
	}

	request.Value["_id"] = result.UpsertedID

	snapshot := bson.M{
		"_uid":       request.Uid,
		"_operation": request.Operation,
		"value":      request.Value,
	}
	publish(sender, snapshot)
}

// Starts watching (observing) a change stream
func _watch(sender interface{}, request api.DocumentRequest) {

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
		log.Fatal("💩 [Unable to watch]: ", err)
	}

	streamContext, _ := context.WithCancel(context.Background())
	go _watchChangeStream(sender, request, streamContext, collectionStream)
}

func _watchChangeStream(sender interface{}, request api.DocumentRequest, context context.Context, stream *mongo.ChangeStream) {
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
			case api.Update, api.Replace:
				request.Query["_id"] = key["_id"].(primitive.ObjectID).Hex()
				_findOne(sender, request)
				return
			case api.Delete:
				doc = bson.M{
					"_id": key["_id"],
				}
				break
			default:
				break
			}
		}

		snapshot := bson.M{
			"_uid":       request.Uid,
			"_operation": request.Operation,
			"value":      doc,
		}
		publish(sender, snapshot)
	}
}
