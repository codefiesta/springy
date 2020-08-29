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
                let message = JSON.parse(e.data);
                self.#broadcast(message);
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
        this.collections.forEach((value, key) => {
            value.notify(message);
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

    // Watches a database reference
    watch = (eventType, callback) => {

        let message = {
            collection: this.name,
            action: SpringyActions.watch,
            operation: eventType
        };

        if (callback) {
            this.subscribers.set(callback, eventType);
        }

        const db = this.database;
        if (db) {
            let payload = JSON.stringify(message);
            db.publish(payload);
        }
    };

    notify = (message) => {

        let key = message['documentKey']['_id'];
        let eventType = message['operationType'];
        this.subscribers.forEach((type, subscriber) => {
            if (type === eventType) {
                subscriber(key, message);
            }
        });
    }

    // Add a new document to this collection with the specified data, assigning it a document ID automatically.
    add = (value, onComplete) => {
        let payload = JSON.stringify(value);
        let message = {
            collection: this.name,
            action: SpringyActions.write,
            operation: SpringyEvents.insert
        };
    };

    find = (onComplete) => {
        
    }
}