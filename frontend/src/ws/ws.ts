import * as ui from "ui";

import * as base from "./base";
import * as types from "./types";
import * as alerts from "../alerts/";

export class WS extends base.BaseWS {
    public events: ui.Events<{
        server: types.WSServer | null;
        open: null;
        close: null;
        message: any;
        "message-devices": types.WSDevice[];
        "message-error": string;
        "message-device": types.WSDevice;
    }> = new ui.Events();

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

    public async request<T extends keyof types.WSCommand>(command: T, data?: types.WSCommand[T]) {
        if (!this.isOpen()) return;
        console.debug(`[ws] Send command: "${command}"`, {
            server: this.server,
            data,
        });

        let request: types.WSRequest = {
            command: command,
            data: data === undefined ? undefined : JSON.stringify(data),
        };

        switch (command) {
            case "GET api.devices":
            case "POST api.device":
            case "PUT api.device":
            case "DELETE api.device":
            case "POST api.device.pins":
            case "POST api.device.color":
                this.ws?.send(JSON.stringify(request));
                break;

            default:
                throw new Error(`unknown command ${command}`);
        }
    }

    protected async handleMessageEvent(ev: MessageEvent) {
        super.handleMessageEvent(ev);
        console.debug("[ws] message.event:", ev);

        if (typeof ev.data === "string") {
            try {
                const resp = JSON.parse(ev.data) as types.WSResponse;
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
                const message = `[ws] Parsing JSON: ${err}`;
                console.warn(message);
                alerts.add("warning", message);
            }
        }

        this.events.dispatch("message", ev.data);
    }

    protected async handleOpenEvent() {
        await super.handleOpenEvent();
        this.events.dispatch("open", null);
    }

    protected async handleCloseEvent() {
        await super.handleCloseEvent();
        this.events.dispatch("close", null);
    }
}
