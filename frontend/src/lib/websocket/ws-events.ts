import { Events } from "ui";
import { BaseWebSocketEvents, type WSServer } from "./base-web-socket-events";

export type WSEvents_Command = {
    "GET api.devices": {
        request: null;
        response: WSEvents_Device[];
    };
    "GET api.device": {
        request: null;
        response: WSEvents_Device;
    };
};

export interface WSEvents_Device {
    server: {
        name: string;
        addr: string;
        isOffline?: boolean;
    };
    pins?: number[] | null;
    color?: number[] | null;
}

export interface WSEvents_Request<T extends keyof WSEvents_Command> {
    command: string;
    data: WSEvents_Command[T]["request"];
}

export interface WSEvents_Response {
    data: any;
    type: "devices" | "device";
}

export class WSEvents extends BaseWebSocketEvents {
    events: Events<{
        server: WSServer | null;
        open: null;
        close: null;
        message: any;
        "message-devices": WSEvents_Command["GET api.devices"]["response"];
        "message-device": WSEvents_Command["GET api.device"]["response"];
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

        switch (command) {
            case "GET api.devices":
                const request: WSEvents_Request<"GET api.devices"> = {
                    command: command,
                    data: null,
                };
                this.ws.send(JSON.stringify(request));
                break;
            default:
                throw new Error(`unknown path ${command}`);
        }
    }

    async handleMessageEvent(ev: MessageEvent) {
        super.handleMessageEvent(ev);
        console.debug("[ws] message.event:", ev);

        if (typeof ev.data === "string") {
            try {
                const req = JSON.parse(ev.data) as WSEvents_Response;
                if (["devices", "device"].includes(req.type)) {
                    this.events.dispatch(`message-${req.type}`, req.data);
                    return;
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
