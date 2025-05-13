export function create(): WS {
    let timeout: NodeJS.Timeout | null = null;
    const timeoutDuration = 1000;

    function getURL(): string {
        return process.env.SERVER_PATH_PREFIX + `/ws`;
    }

    const onClose = () => {
        if (timeout !== null) {
            clearTimeout(timeout);
            timeout = null;
        }

        window.utils.setOnlineIndicatorState(false);

        // Reconnect here
        timeout = setTimeout(async () => {
            console.debug(`Try to reconnect to "${getURL()}"`);
            await ws.connect();
        }, timeoutDuration);
    };

    const onOpen = () => {
        console.debug(`WS: connected to "${getURL()}"`);
        window.utils.setOnlineIndicatorState(true);
        ws.events.dispatch("open", undefined);
    };

    const onMessage = async (ev: MessageEvent<Blob>) => {
        const data: WSMessageData = JSON.parse(await ev.data.text());
        console.debug(`WS: Got a message:`, data);

        switch (data.type) {
            case "device":
                {
                    window.store.obj.update("devices", (devices) => {
                        // Update device in store
                        for (let x = 0; x < devices.length; x++) {
                            if (
                                devices[x].server.addr !== data.data.server.addr
                            ) {
                                continue;
                            }

                            devices[x] = data.data;
                        }

                        return devices;
                    });
                }
                break;
            case "colors":
                {
                    window.store.obj.set("colors", data.data);
                }
                break;
        }

        ws.events.dispatch(data.type, data.data);
    };

    const ws: WS = {
        events: new window.ui.Events(),
        socket: null,

        isOpen() {
            if (!this.socket) return false;
            return this.socket.readyState === this.socket.OPEN;
        },

        async connect() {
            if (this.socket) this.close();

            const wsAddr = getURL(); // origin + path
            console.debug(`Try to connect WebSocket to ${wsAddr}`);

            this.socket = new WebSocket(wsAddr);

            // Reconnect handler
            this.socket.addEventListener("close", onClose);
            this.socket.addEventListener("open", onOpen);
            this.socket.addEventListener("message", onMessage);
        },

        close() {
            if (timeout) {
                clearTimeout(timeout);
                timeout = null;
            }

            if (this.socket) {
                this.socket.removeEventListener("close", onClose);
                if (this.isOpen()) this.socket.close();
                this.socket = null;
            }
        },
    };

    return ws;
}
