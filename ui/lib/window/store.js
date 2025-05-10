/**
 * @returns {import("../../types").Store}
 */
export function create() {
    /** @type {import("../../types").UIStore} */
    const store = new window.ui.Store("picow-led:");

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
                        data[device.server.addr] = device.color;
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

        /**
         * @param {string} addr
         * @returns {import("../../types").Device | null}
         */
        device(addr) {
            for (const d of store.get("devices") || []) {
                if (d.server.addr === addr) {
                    return d;
                }
            }

            return null;
        },

        /**
         * @param {string} addr
         * @returns {import("../../types").Color | null}
         */
        currentDeviceColor(addr) {
            const currentDeviceColors = store.get("currentDeviceColors");
            for (const key in currentDeviceColors) {
                if (key === addr) {
                    return currentDeviceColors[key];
                }
            }

            return null;
        },
    };
}
