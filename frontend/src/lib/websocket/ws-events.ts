import { Events } from "ui";
import { BaseWebSocketEvents } from "./base-web-socket-events";
import type {
    WSEvents_Command,
    WSEvents_Device,
    WSEvents_Request,
    WSEvents_Response,
    WSEvents_Server,
} from "./types";

export class WSEvents extends BaseWebSocketEvents {
    events: Events<{
        server: WSEvents_Server | null;
        open: null;
        close: null;
        message: any;
        "message-devices": WSEvents_Device[];
        "message-error": string;
        "message-device": WSEvents_Device;
    }> = new Events();

    constructor() {
        super("/ws");
    }

    get server() {
        return super.server;
    }

    set server(value) {
        super.server = value;
        this.events.dispatch("server", value);
    }

    async request<T extends keyof WSEvents_Command>(
        command: T,
        data: WSEvents_Command[T]["request"] = null
    ) {
        if (!this.isOpen()) return;
        console.debug(`[ws] Send command: "${command}"`, {
            server: this.server,
            data,
        });

        let request: WSEvents_Request;
        switch (command) {
            case "GET api.devices":
                request = {
                    command: command,
                    data: null,
                };
                this.ws.send(JSON.stringify(request));
                break;
            case "POST api.device.color":
                request = {
                    command: command,
                    data: JSON.stringify(data),
                };
                this.ws.send(JSON.stringify(request));
                break;
            default:
                throw new Error(`unknown command ${command}`);
        }
    }

    async handleMessageEvent(ev: MessageEvent) {
        super.handleMessageEvent(ev);
        console.debug("[ws] message.event:", ev);

        if (typeof ev.data === "string") {
            try {
                const resp = JSON.parse(ev.data) as WSEvents_Response;
                console.debug(`[ws] message:`, resp);

                switch (resp.type) {
                    case "devices":
                    case "device":
                        this.events.dispatch(`message-${resp.type}`, resp.data);
                        break;
                    case "error":
                        this.events.dispatch(`message-error`, resp.data);
                        break;
                }
            } catch (err) {
                console.warn("[ws] Parsing JSON:", err);
            }
        }

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
