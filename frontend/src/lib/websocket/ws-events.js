import { Events } from "ui";
import { BaseWebSocketEvents } from "./base-web-socket-events";

export class WSEvents extends BaseWebSocketEvents {
    constructor() {
        super("/ws");

        /**
         * @type {Events<{
         *  "server": Server | null;
         *  "open": null;
         *  "close": null;
         *  "message": any;
         *  "message-device": Device;
         *  "message-devices": Device[];
         *  }>}
         */
        this.events = new Events();
        // TODO: Dispatch "message-device" and "message-devices" events
    }

    get server() {
        return super.server;
    }

    set server(value) {
        super.server = value;
        this.events.dispatch("server", value);
    }

    /** @param {MessageEvent} ev */
    async handleMessageEvent(ev) {
        super.handleMessageEvent(ev);
        // TODO: Parsing data and dispatch "message-device" or "message-devices"

        //if (ev.data instanceof Blob) {
        //    this.ws.send("pong");
        //    return;
        //}

        //const device = JSON.parse(ev.data);
        //this.events.dispatch("message", device);
        this.events.dispatch("message", ev.data);
    }

    async handleOpenEvent() {
        await super.handleOpenEvent();
        this.events.dispatch("open", null);
    }

    async handleCloseEvent() {
        await super.handleCloseEvent();
        this.events.dispatch("close", null);
    }
}
