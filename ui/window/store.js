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

    return {
        obj: store,

        /**
         * @param {string} addr
         */
        device(addr) {
            for (const d of store.get("devices") || []) {
                if (d.server.addr === addr) {
                    return d;
                }
            }

            throw new Error(`device ${addr} not found`);
        },
    };
}
