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
         * @returns {import("../types.d.ts").Device | null}
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
         * @returns {import("../types.d.ts").Color | null}
         */
        currentColor(addr) {
            const current = store.get("color").current;
            for (const a in current) {
                if (a === addr) {
                    return current[a];
                }
            }

            return null;
        },
    };
}
