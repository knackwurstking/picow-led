import { WSServer } from "./types";

export class BaseWS {
    private _server: WSServer | null = null;

    private closeHandler = async () => {
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
        return this._server;
    }

    set server(value) {
        if (!value) {
            this.origin = "";
        } else {
            const addr = !value.port ? value.host : `${value.host}:${value.port}`;

            this.origin = `${value.ssl ? "wss:" : "ws:"}//${addr}`;
        }

        this.connect();
        this._server = value;
    }

    protected isOpen() {
        if (!this.ws) return false;
        return this.ws.readyState === this.ws.OPEN;
    }

    protected connect() {
        if (this.ws) this.close();
        const addr = this.origin + this.path;
        console.debug(`Connect WebSocket to ${addr}`);
        this.ws = new WebSocket(addr);

        this.ws.addEventListener("message", (ev) => {
            this.handleMessageEvent(ev);
        });

        this.ws.addEventListener("open", () => {
            this.handleOpenEvent();
        });

        this.ws.addEventListener("error", (ev) => {
            this.handleErrorEvent(ev);
        });

        this.ws.addEventListener("close", this.closeHandler);
    }

    protected close() {
        this.ws?.removeEventListener("close", this.closeHandler);

        if (!!this.timeout) {
            clearTimeout(this.timeout);
            this.timeout = null;
        }

        if (this.isOpen()) this.ws?.close();
    }

    protected async handleMessageEvent(_ev: MessageEvent) {}

    protected async handleOpenEvent() {
        console.debug(`websocket connection established "${this.origin}${this.path}"`);
    }

    protected async handleErrorEvent(ev: Event) {
        console.error(`websocket connection error "${this.origin}${this.path}"`, ev);
    }

    protected async handleCloseEvent() {
        console.warn(`websocket connection closed "${this.origin}${this.path}"`);

        this.timeout = setTimeout(async () => {
            this.connect();
        }, this.timeoutDuration);
    }
}
