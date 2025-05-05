/**
 * @returns {import("../types.d.ts").Store}
 */
export function create() {
    /** @type {import("../types.d.ts").PageWindow} */
    // @ts-expect-error
    const w = window;

    /** @type {import("../types.d.ts").UIStore} */
    const store = new w.ui.Store("picow-led:");

    store.set("devices", [], true);
    store.set(
        "color",
        {
            api: [],
            current: {},
        },
        true,
    );

    return {
        obj: store,

        /**
         * @param {string} addr
         * @returns {import("../types.d.ts").Device}
         */
        device(addr) {
            for (const d of store.get("devices") || []) {
                if (d.server.addr === addr) {
                    return d;
                }
            }

            throw new Error(`device ${addr} not found`);
        },

        /**
         * @param {string} addr
         * @returns {import("../types.d.ts").Color}
         */
        currentColor(addr) {
            const current = store.get("color").current;
            for (const a in current) {
                if (a === addr) {
                    return current[a];
                }
            }

            throw new Error(`current color for ${addr} not found`);
        },
    };
}
