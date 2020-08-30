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
                console.error(err);
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

    // Watches a collection for events
    watch = (eventType, callback) => {

        let subscriber = new DataSubscriber(this.name, SpringyActions.watch, eventType, callback);
        this.subscribers.set(subscriber.identifier, subscriber);

        let encoded = subscriber.encode();
        this.database.publish(encoded);
        return subscriber;
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

    // Add a new document to this collection with the specified data, assigning it a document ID automatically.
    add = (value, callback) => {

        let subscriber = new DataSubscriber(this.name, SpringyActions.write, SpringyEvents.insert, callback);
        this.subscribers.set(subscriber.identifier, subscriber);

        let encoded = subscriber.encode(value);
        this.database.publish(encoded);
    };

    find = (callback) => {

    }
}

class DataSubscriber {

    constructor(collection, action, event, callback) {
        this.sid = uuidv4();
        this.collection = collection;
        this.action = action;
        this.event = event;
        this.callback = callback;
    }

    encode = (value) => {
        let encoded = {
            _sid: this.sid,
            collection: this.collection,
            action: this.action,
            operation: this.event,
            value: value
        };
        return JSON.stringify(encoded);
    }

    get identifier() {
        return this.sid;
    }
}

class DataSnapshot {

    constructor(data) {
        this.sid = data["_sid"];
        this.key = data["key"];
        this.value = data["value"] ?? {};
    }

    get identifier() {
        return this.sid;
    }
}

function uuidv4() {
    return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
}