
class Springy {

    constructor(config) {

        this.callbacks = [];
        this.isConnected = false;

        if (!window["WebSocket"]) {
            console.log('💩 Your browser does not support WebSockets.');
            return;
        }

        this.ws = new WebSocket(config.databaseURL);
        this.#addHandlers();
    }

    #addHandlers = () => {

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
                self.#routeMessage(message);
            } catch(err) {
                console.error(err);
            }
        };
    };

    watch = (path, eventType, callback) => {

        let message = {
            path: path,
            action: "watch",
            operation: eventType
        };

        if (callback) { this.callbacks.push(callback); }
        let payload = JSON.stringify(message);
        this.#send(payload);
    };

    /// Internal Message Handling

    #routeMessage = (message) => {
        // Handle the incoming message
        for (const callback of this.callbacks) {
            callback(message);
        }
    };

    #send = (message) => {
        this.#queueMessage(() => {
            this.ws.send(message);
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
}