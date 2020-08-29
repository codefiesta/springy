
class Springy {

    constructor(config) {
        if (!window["WebSocket"]) {
            console.log('💩 Your browser does not support WebSockets.');
            return;
        }
        this._database = new Database(config);
    }

    // Returns our enum of events
    static get events() {
        return Object.freeze({
            insert: "insert",
            update: "update",
            remove: "delete",
            replace: "replace",
        });
    }

    // Returns the app database object
    get database() {
        return this._database;
    }
}

class Database {

    constructor(config) {
        this.isConnected = false;
        this.references = new Map();
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
            } catch(err) {
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

    /// Broadcasts an incoming message to reference handlers
    #broadcast = (message) => {
        this.references.forEach((value, key) => {
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

    // Returns a data reference for the specified path
    reference = (path) => {
        if (this.references.has(path)) {
            return this.references.get(path);
        }
        let reference = new DataReference(this, path);
        this.references.set(path, reference);
        return reference;
    };
}

class DataReference {

    constructor(database, path) {
        this.database = database;
        this.subscribers = [];
        this.path = path;
    }

    child = (path) => {
        return new DataReference(this.database, path);
    }

    // Watches a database reference
    watch = (eventType, callback) => {

        let message = {
            path: this.path,
            action: "watch",
            operation: eventType
        };

        console.log(`👀 ${eventType}`);
        if (callback) { this.subscribers.push(callback); }

        const db = this.database;
        if (db) {
            let payload = JSON.stringify(message);
            db.publish(payload);
        }
    };

    notify = (message) => {
        for (const subscriber of this.subscribers) {
            subscriber(message);
        }
    }

    // Writes data to this Database location.
    // This will overwrite any data at this location and all child locations.
    // Passing null for the new value is equivalent to calling remove()
    set = (value, onComplete) => {

    };
}