export declare type WSMessageData =
    | {
          type: "device";
          data: Device;
      }
    | {
          type: "colors";
          data: Colors;
      };

export class WS {
    public socket: WebSocket | null = null;
    public events = new window.ui.Events<{
        open: undefined;
        close: undefined;
        device: Device;
        colors: Colors;
    }>();

    protected timeout: NodeJS.Timeout | null = null;
    protected timeoutDuration = 1000;
    protected open = false;

    protected onClose = () => {
        if (this.timeout !== null) {
            clearTimeout(this.timeout);
            this.timeout = null;
        }

        window.utils.setOnlineIndicatorState(false);

        if (this.open) {
            this.open = false;
            this.events.dispatch("close", undefined);
        }

        // Reconnect here
        this.timeout = setTimeout(async () => {
            console.debug(`Try to reconnect to "${this.getURL()}"`);
            await this.connect();
        }, this.timeoutDuration);
    };

    protected onOpen = () => {
        console.debug(`WS: connected to "${this.getURL()}"`);
        window.utils.setOnlineIndicatorState(true);

        if (!this.open) {
            this.open = true;
            this.events.dispatch("open", undefined);
        }
    };

    protected onMessage = async (ev: MessageEvent<Blob>) => {
        const data: WSMessageData = JSON.parse(await ev.data.text());
        console.debug(`WS: Got a message:`, data);

        switch (data.type) {
            case "device":
                await this.wsHandleDevice(data.data);
                break;
            case "colors":
                await this.wsHandleColors(data.data);
                break;
        }

        this.events.dispatch(data.type, data.data);
    };

    protected getURL(): string {
        return process.env.SERVER_PATH_PREFIX + `/ws`;
    }

    public isOpen() {
        return this.open;
    }

    public async connect() {
        if (this.socket) this.close();

        const wsAddr = this.getURL(); // origin + path
        console.debug(`Try to connect WebSocket to ${wsAddr}`);

        this.socket = new WebSocket(wsAddr);

        // Reconnect handler
        this.socket.addEventListener("close", this.onClose);
        this.socket.addEventListener("open", this.onOpen);
        this.socket.addEventListener("message", this.onMessage);
    }

    public close() {
        if (this.timeout) {
            clearTimeout(this.timeout);
            this.timeout = null;
        }

        if (this.socket) {
            this.socket.removeEventListener("close", this.onClose);
            if (this.isOpen()) this.socket.close();
            this.socket = null;
        }
    }

    protected async wsHandleDevice(data: Device): Promise<void> {
        window.store.update("devices", (devices) => {
            // Update device in store
            for (let x = 0; x < devices.length; x++) {
                if (devices[x].server.addr !== data.server.addr) {
                    continue;
                }

                devices[x] = data;
            }

            return devices;
        });
    }

    protected async wsHandleColors(data: Colors): Promise<void> {
        window.store.set("colors", data);
    }
}
