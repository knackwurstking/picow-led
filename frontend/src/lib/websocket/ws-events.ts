import { Events } from "ui";
import { BaseWebSocketEvents } from "./base-web-socket-events";

type WSEvents_Command = "GET api.devices";

export class WSEvents extends BaseWebSocketEvents {
    events: Events<{
        server: Server | null;
        open: null;
        close: null;
        message: any;
        messageDevice: Device;
        messageDevices: Device[];
    }>;

    constructor() {
        super("/ws");
        this.events = new Events();
    }

    get server() {
        return super.server;
    }

    set server(value) {
        super.server = value;
        this.events.dispatch("server", value);
    }

    async request(command: WSEvents_Command, data: any = null) {
        if (!this.isOpen()) return;
        console.debug(`[ws] Send command: "GET api.devices"`, this.server);

        switch (command) {
            case "GET api.devices":
                // TODO: Need request type data here `{ command: string; data: any }`
                this.ws.send(`GET api.devices`);
                break;
            default:
                throw new Error(`unknown path ${command}`);
        }
    }

    async handleMessageEvent(ev: MessageEvent) {
        super.handleMessageEvent(ev);
        console.debug("[ws] message.event:", ev);

        // TODO: Parsing data and dispatch "message-device" or "message-devices"

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
