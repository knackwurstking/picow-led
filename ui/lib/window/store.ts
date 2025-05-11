export function create(): Store {
    const store: UIStore = new window.ui.Store("picow-led:");

    store.set("devices", [], true);
    store.set("colors", [], true);

    store.set("currentDeviceColors", {}, true);

    // Update "currentDeviceColors" after any "devices" change
    store.listen(
        "devices",
        async (devices) => {
            devices.forEach((device) => {
                if (Math.max(...(device.color || [])) > 0) {
                    store.update("currentDeviceColors", (data) => {
                        data[device.server.addr] =
                            device.color || (device.pins || []).map(() => 255);
                        return data;
                    });
                }
            });
        },
        true,
    );

    // @ts-expect-error
    store.delete("color"); // TODO: Just to clean up, can be removed before release

    return {
        obj: store,

        // Helper Functions:

        device(addr: string): Device | null {
            for (const d of store.get("devices") || []) {
                if (d.server.addr === addr) {
                    return d;
                }
            }

            return null;
        },

        currentDeviceColor(addr: string): Color | null {
            return store.get("currentDeviceColors")?.[addr] || null;
        },
    };
}
