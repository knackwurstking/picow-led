import type { WSEventsServer } from "./types";

export class BaseWebSocketEvents {
    #server: WSEventsServer | null = null;

    #messageHandler = async (ev: MessageEvent) => {
        await this.handleMessageEvent(ev);
    };
    #openHandler = async () => {
        await this.handleOpenEvent();
    };
    #errorHandler = async (ev: Event) => {
        await this.handleErrorEvent(ev);
    };
    #closeHandler = async () => {
        await this.handleCloseEvent();
    };

    path: string;
    origin: string;

    timeout: NodeJS.Timeout | null;
    timeoutDuration: number;

    ws: WebSocket | null;

    constructor(path: string) {
        this.path = path;
        this.origin = "";

        this.timeout = null;
        this.timeoutDuration = 1000;

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

        if (!!this.timeout) {
            clearTimeout(this.timeout);
            this.timeout = null;
        }

        if (this.isOpen()) this.ws.close();
    }

    async handleMessageEvent(ev: MessageEvent) {}

    async handleOpenEvent() {
        console.debug(
            `websocket connection established "${this.origin}${this.path}"`
        );
    }

    async handleErrorEvent(ev: Event) {
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
