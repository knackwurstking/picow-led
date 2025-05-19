window.addEventListener("load", () => {
    let cleanup: import("ui").CleanUpFunction[] = [];

    window.addEventListener("pageshow", () => {
        // Starting the WebSocket handler
        window.ws.connect();

        cleanup.push(
            window.ws.events.addListener("open", () => {
                window.utils.setOnlineIndicatorState(true);
            }),

            window.ws.events.addListener("close", () => {
                window.utils.setOnlineIndicatorState(false);
            }),

            window.ws.events.addListener("message", async (data) => {
                switch (data.type) {
                    case "device":
                        await wsHandleDevice(data.data);
                        break;
                    case "colors":
                        await wsHandleColors(data.data);
                        break;
                }
            }),
        );
    });

    window.addEventListener("pagehide", () => {
        cleanup.forEach((fn) => fn());
        cleanup = [];
    });
});

async function wsHandleDevice(data: Device): Promise<void> {
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

async function wsHandleColors(data: Colors): Promise<void> {
    window.store.set("colors", data);
}
