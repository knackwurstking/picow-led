import { Events } from "ui";
import { BaseWebSocketEvents, type WSServer } from "./base-web-socket-events";

export type WSEvents_Command = {
    "GET api.devices": {
        request: null;
        response: WSEvents_Device[];
    };
    "POST api.device.color": {
        request: {
            addr: string;
            color: number[];
        };
        response: null;
    };
};

export interface WSEvents_DeviceServer {
    name?: string;
    addr: string;
    online?: boolean;
}

export interface WSEvents_Device {
    server: WSEvents_DeviceServer;
    pins?: number[];
    color?: number[];
}

export interface WSEvents_Request {
    command: string;
    data: string; // NOTE: JSON string
}

export type WSEvents_Response =
    | {
          data: string;
          type: "error";
      }
    | {
          data: WSEvents_Device[];
          type: "devices";
      }
    | {
          data: WSEvents_Device;
          type: "device";
      };

export class WSEvents extends BaseWebSocketEvents {
    events: Events<{
        server: WSServer | null;
        open: null;
        close: null;
        message: any;
        "message-devices": WSEvents_Device[];
        "message-error": string;
        "message-device": WSEvents_Device;
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

    async request<T extends keyof WSEvents_Command>(
        command: T,
        data: WSEvents_Command[T]["request"] = null
    ) {
        if (!this.isOpen()) return;
        console.debug(`[ws] Send command: "GET api.devices"`, this.server);

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
