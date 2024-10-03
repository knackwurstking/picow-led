export class BaseWebSocketEvents {
    /**
     * @type {Server | null}
     */
    #server = null;

    /**
     * @param {MessageEvent} ev
     */
    #messageHandler = async (ev) => {
        await this.handleMessageEvent(ev);
    };

    #openHandler = async () => {
        await this.handleOpenEvent();
    };

    /**
     * @param {Event} ev
     */
    #errorHandler = async (ev) => {
        await this.handleErrorEvent(ev);
    };

    /**
     * @type {() => Promise<void>}
     */
    #closeHandler = async () => {
        await this.handleCloseEvent();
    };

    /**
     * @param {string} path
     */
    constructor(path) {
        this.path = path;
        this.origin = "";

        /** @type {NodeJS.Timeout} */
        this.timeout;
        this.timeoutDuration = 1000;

        /** @type {WebSocket | null} */
        this.ws = null;
    }

    get server() {
        return this.#server;
    }

    set server(value) {
        if (!value) {
            this.origin = "";
        } else {
            const addr = !value.port
                ? value.host
                : `${value.host}:${value.port}`;

            this.origin = `${value.ssl ? "wss:" : "ws:"}//${addr}`;
        }

        this.connect();
        this.#server = value;
    }

    isOpen() {
        if (!this.ws) return false;
        return this.ws.readyState === this.ws.OPEN;
    }

    connect() {
        if (this.ws) this.close();
        this.ws = new WebSocket(this.origin + this.path);

        this.ws.addEventListener("message", this.#messageHandler);
        this.ws.addEventListener("open", this.#openHandler);
        this.ws.addEventListener("error", this.#errorHandler);
        this.ws.addEventListener("close", this.#closeHandler);
    }

    close() {
        this.ws.removeEventListener("close", this.#closeHandler);
        if (!!this.timeout) clearTimeout(this.timeout);
        if (this.isOpen()) this.ws.close();
    }

    /**
     * @param {MessageEvent} ev
     */
    async handleMessageEvent(ev) {}

    async handleOpenEvent() {
        console.debug(
            `websocket connection established "${this.origin}${this.path}"`
        );
    }

    /**
     * @param {Event} ev
     */
    async handleErrorEvent(ev) {
        console.error(
            `websocket connection error "${this.origin}${this.path}"`,
            ev
        );
    }

    async handleCloseEvent() {
        console.warn(
            `websocket connection closed "${this.origin}${this.path}"`
        );

        this.timeout = setTimeout(async () => {
            this.connect();
        }, this.timeoutDuration);
    }
}
