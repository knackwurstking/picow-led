import { Events } from "ui";
import { WebSocketEvents } from "./web-socket-events";

export class DevicesEvents extends WebSocketEvents {
    constructor() {
        super("/events/devices");

        /**
         * @type {Events<{
         *  "server": Server | null;
         *  "open": null;
         *  "close": null;
         *  "message": Device[]
         *  }>}
         */
        this.events = new Events();
    }

    get server() {
        return super.server;
    }

    set server(value) {
        super.server = value;
        this.events.dispatch("server", value);
    }

    /** @param {MessageEvent} ev */
    async onMessage(ev) {
        super.onMessage(ev);
        if (ev.data instanceof Blob) {
            this.ws.send("pong");
            return;
        }

        const devices = JSON.parse(ev.data);
        this.events.dispatch("message", devices);
    }

    async onOpen() {
        await super.onOpen();
        this.events.dispatch("open", null);
    }

    async onClose() {
        await super.onClose();
        this.events.dispatch("close", null);
    }
}
