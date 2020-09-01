const SpringyScope = Object.freeze({
    find: "find",
    findOne: "findOne",
    write: "write",
    watch: "watch",
});

const SpringyEvents = Object.freeze({
    insert: "insert",
    update: "update",
    delete: "delete",
    replace: "replace",
});

class Springy {

    constructor(config) {
        if (!window["WebSocket"]) {
            console.log('💩 Your browser does not support WebSockets.');
            return;
        }
        this._database = new Database(config);
    }

    // Returns the app database object
    get database() {
        return this._database;
    }
}

class Database {

    constructor(config) {
        this.isConnected = false;
        this.collections = new Map();
        this.ws = new WebSocket(config.databaseURL);
        this.addSocketHandlers();
    }

    addSocketHandlers = () => {

        let self = this;
        this.ws.onopen = function (e) {
            self.isConnected = true;
        };
        this.ws.onclose = function (e) {
            self.isConnected = false;
        };
        this.ws.onmessage = function (e) {
            try {
                let data = JSON.parse(e.data);
                data.forEach(message => {
                    self.broadcast(message);
                });
            } catch (err) {
                console.error(err, e.data);
            }
        };
    };

    /// Publishes a message out to the database
    publish = (message) => {
        this.queueMessage(() => {
            this.ws.send(message);
        });
    };

    /// Broadcasts an incoming message to collection handlers
    broadcast = (message) => {
        this.collections.forEach((collection, key) => {
            collection.notify(message);
        });
    };

    queueMessage = (callback) => {
        if (this.ws.readyState === 1) {
            callback();
        } else {
            let self = this;
            setTimeout(() => {
                self.queueMessage(callback);
            }, 5);
        }
    };

    // Returns a database collection for the specified name
    collection = (name) => {
        if (this.collections.has(name)) {
            return this.collections.get(name);
        }
        let collection = new DocumentCollection(this, name);
        this.collections.set(name, collection);
        return collection;
    };
}

class DocumentCollection {

    constructor(database, name) {
        this.database = database;
        this.subscribers = new Map();
        this.name = name;
    }

    // Queues the subscriber event
    subscribe = (subscriber) => {
        this.subscribers.set(subscriber.identifier, subscriber);
        let encoded = subscriber.encode();
        this.database.publish(encoded);
    }

    // Watches a collection for events
    watch = (eventType, callback) => {
        let subscriber = new DataSubscriber(this.name, {}, SpringyScope.watch, eventType, null, callback);
        this.subscribe(subscriber);
        return subscriber;
    };

    // Fetches all documents in the collection
    get = (callback) => {
        let subscriber = new DataSubscriber(this.name, {}, SpringyScope.find, null, null, callback);
        this.subscribe(subscriber);
    };


    // Notifies all interested subscribers that we received a collection event
    notify = (data) => {
        let snapshot = new DataSnapshot(this, data);
        if (this.subscribers.has(snapshot.identifier)) {
            let subscriber = this.subscribers.get(snapshot.identifier);
            if (subscriber.callback) {
                subscriber.callback(snapshot);
            }

            // REMOVE ANY SINGLE FIRE READ OR WRITE SUBSCRIBERS
            switch (subscriber.scope) {
                case SpringyScope.write, SpringyScope.find, SpringyScope.findOne:
                    this.subscribers.delete(snapshot.identifier);
                    break;
                default:
                    break;
            }
        }
    }

    // Add a new document to this collection with the specified data, assigning it a document ID automatically.
    add = (value, callback) => {
        let subscriber = new DataSubscriber(this.name, {}, SpringyScope.write, SpringyEvents.insert, value, callback);
        this.subscribe(subscriber);
    }

    // Removes a document with the specified key
    remove = (key, callback) => {
        let query = {"_id": key};
        let subscriber = new DataSubscriber(this.name, query, SpringyScope.write, SpringyEvents.delete, null, callback);
        this.subscribe(subscriber);
    }
}

class DataSubscriber {

    constructor(collection, query, scope, event, value, callback, onDisconnect) {
        this.uid = uuidv4();
        this.collection = collection;
        this.query = query;
        this.scope = scope;
        this.event = event;
        this.value = value;
        this.onDisconnect = onDisconnect ?? false;
        this.callback = callback;
    }

    encode = () => {
        let encoded = {
            _uid: this.uid,
            collection: this.collection,
            query: this.query,
            scope: this.scope,
            operation: this.event,
            value: this.value ?? {},
            onDisconnect: this.onDisconnect
        };
        return JSON.stringify(encoded);
    }

    get identifier() {
        return this.uid;
    }
}

// Contains data read from a document in the database
class DataSnapshot {

    constructor(collection, data) {
        this.collection = collection;
        this.uid = data["_uid"];
        this.value = data["value"] ?? {};
        this._onDisconnect = new OnDisconnect(this);
    }

    get key() {
        return this.value._id;
    }

    get identifier() {
        return this.uid;
    }

    onDisconnect = () => {
        return this._onDisconnect;
    }
}


class OnDisconnect {

    constructor(snapshot) {
        this.snapshot = snapshot;
    }

    remove = () => {
        // Generate a message to send to the database
        let query = {"_id": this.snapshot.key};
        let subscriber = new DataSubscriber(this.snapshot.collection.name, query, SpringyScope.write, SpringyEvents.delete, null, null, true);
        this.snapshot.collection.subscribe(subscriber);
    }

    set = (value) => {
        let query = {"_id": this.snapshot.key};
        let subscriber = new DataSubscriber(this.snapshot.collection.name, query, SpringyScope.write, SpringyEvents.update, value, null, true);
        this.snapshot.collection.subscribe(subscriber);
    }
}

function uuidv4() {
    return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
}