const SpringyEvents = Object.freeze({
    insert: "insert",
    update: "update",
    delete: "delete",
    replace: "replace",
});

const SpringyActions = Object.freeze({
    read: "read",
    write: "write",
    watch: "watch",
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
        this.#addSocketHandlers();
    }

    #addSocketHandlers = () => {

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
                self.#broadcast(data);
            } catch (err) {
                console.error(err, e.data);
            }
        };
    };

    /// Publishes a message out to the database
    publish = (message) => {
        this.#queueMessage(() => {
            this.ws.send(message);
        });
    };

    /// Broadcasts an incoming message to collection handlers
    #broadcast = (message) => {
        this.collections.forEach((collection, key) => {
            collection.notify(message);
        });
    };

    #queueMessage = (callback) => {
        if (this.ws.readyState === 1) {
            callback();
        } else {
            let self = this;
            setTimeout(() => {
                self.#queueMessage(callback);
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
        let subscriber = new DataSubscriber(this.name, null, SpringyActions.watch, eventType, null, callback);
        this.subscribe(subscriber);
        return subscriber;
    };

    // Fetches all documents in the collection
    get = (callback) => {
        let subscriber = new DataSubscriber(this.name, null, SpringyActions.read, null, null, callback);
        this.subscribe(subscriber);
    };


    // Notifies all interested subscribers that we received a collection event
    notify = (data) => {
        let snapshot = new DataSnapshot(data);
        if (this.subscribers.has(snapshot.identifier)) {
            let subscriber = this.subscribers.get(snapshot.identifier);
            if (subscriber.callback) {
                subscriber.callback(snapshot);
            }

            // REMOVE ANY SINGLE FIRE READ OR WRITE SUBSCRIBERS
            switch (subscriber.action) {
                case SpringyActions.write:
                    this.subscribers.delete(snapshot.identifier);
                    break;
                case SpringyActions.read:
                    this.subscribers.delete(snapshot.identifier);
                    break;
                default:
                    break;
            }
        }
    }

    // Returns a document in this collection with the specified identifier.
    // If no identifier is specified, an automatically-generated unique ID will be used for the returned doc.
    doc = (key) => {
        let document = new Document(this, key);
        return document;
    }

    // Add a new document to this collection with the specified data, assigning it a document ID automatically.
    add = (value, callback) => {
        let subscriber = new DataSubscriber(this.name, null, SpringyActions.write, SpringyEvents.insert, value, callback);
        this.subscribe(subscriber);
    }

    // Removes a document with the specified key
    remove = (key, callback) => {
        let subscriber = new DataSubscriber(this.name, key, SpringyActions.write, SpringyEvents.delete, null, callback);
        this.subscribe(subscriber);
    }
}


class Document {

    constructor(collection, key) {
        this.collection = collection;
        this.key = key;
        this._onDisconnect = new OnDisconnect(this);
    }

    // Writes data to this document location.
    set = (value, callback) => {
        let operation = this.key === undefined ? SpringyEvents.insert : SpringyEvents.update
        let subscriber = new DataSubscriber(this.collection.name, this.key, SpringyActions.write, operation, value, callback);
        this.collection.subscribe(subscriber);
    }

    onDisconnect = () => {
        return this._onDisconnect;
    }
}

class OnDisconnect {

    constructor(doc) {
        this.doc = doc;
    }

    remove = () => {

    }

    set = (value, callback) => {

    }

}

class DataSubscriber {

    constructor(collection, key, action, event, value, callback, onDisconnect) {
        this.uid = uuidv4();
        this.collection = collection;
        this.key = key;
        this.action = action;
        this.event = event;
        this.value = value;
        this.onDisconnect = onDisconnect ?? false;
        this.callback = callback;
    }

    encode = () => {
        let encoded = {
            _uid: this.uid,
            collection: this.collection,
            key: this.key,
            action: this.action,
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

    constructor(data) {
        this.uid = data["_uid"];
        this.key = data["key"];
        this.value = data["value"] ?? {};
    }

    get identifier() {
        return this.uid;
    }
}

function uuidv4() {
    return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
}